import QtQuick
import QtQuick.Controls
import QtQuick.Layouts

ApplicationWindow {
    id: mainWindow
    width: 400
    height: 800
    visible: true
    title: "Kingo VPN"
    color: "#121212"

    onClosing: (close) => {
        close.accepted = false
        hide()
        trayIcon.showMessage("Kingo VPN", "Minimized to tray.")
    }

    property string connStatus: vpnController ? vpnController.status : "loading"
    property int currentTab: 0 // 0: Favorites, 1: Servers, 2: Custom
    property string selectedServerID: ""
    property real prevDownload: 0
    property real prevUpload: 0
    property real prevTime: 0

    readonly property color bgColor: "#121212"
    readonly property color cardColor: "#1E1E1E"
    readonly property color accentColor: "#2196F3"
    readonly property color textColor: "#FFFFFF"
    readonly property color subTextColor: "#B0B0B0"
    readonly property color successColor: "#4CAF50"

    Timer {
        interval: 1000
        running: connStatus == "connected"
        repeat: true
        onTriggered: vpnController.getTraffic()
    }

    function formatSpeed(bytes) {
        if (!bytes || bytes < 0) return "0.0 B/s"
        if (bytes < 1024) return bytes.toFixed(1) + " B/s"
        if (bytes < 1048576) return (bytes / 1024).toFixed(1) + " KB/s"
        return (bytes / 1048576).toFixed(1) + " MB/s"
    }

    function getPingColor(latency) {
        if (latency == 9999 || latency == 0) return subTextColor
        if (latency < 150) return successColor
        if (latency < 300) return "#FFC107"
        return "#F44336"
    }

    function getButtonColor() {
        if (connStatus == "connected") return successColor
        if (connStatus == "connecting") return "#FFC107"
        return accentColor
    }
    
    function getStatusText() {
        if (connStatus == "connecting") return "Connecting..."
        if (connStatus == "connected") return "Connected"
        return "Disconnected"
    }
    
    function getSubStatusText() {
        if (connStatus == "connected") return ""
        return "Tap to connect"
    }

    Connections {
        target: trayIcon
        function onActivateRequested() { mainWindow.show(); mainWindow.raise(); mainWindow.requestActivate() }
        function onConnectRequested() { mainWindow.show(); mainWindow.raise() }
        function onDisconnectRequested() { vpnController.disconnectVpn() }
        function onQuitRequested() { Qt.quit() }
    }

    Connections {
        target: vpnController
        function onTrafficChanged() {
            var currTime = Date.now()
            var dt = (currTime - prevTime) / 1000.0
            if (dt > 0) {
                downSpeedText.text = formatSpeed((vpnController.downloadSpeed - prevDownload) / dt)
                upSpeedText.text = formatSpeed((vpnController.uploadSpeed - prevUpload) / dt)
            }
            prevDownload = vpnController.downloadSpeed
            prevUpload = vpnController.uploadSpeed
            prevTime = currTime
        }
    }

    Menu {
        id: serverMenu
        background: Rectangle { color: cardColor; border.color: "#333333" }
        MenuItem { text: "Add Server"; onTriggered: addServerDialog.open() }
        MenuItem { text: "Fetch Servers"; onTriggered: vpnController.addSubscription("https://raw.githubusercontent.com/MhdiTaheri/VpnHub/main/sub") }
        MenuItem { text: "Test All Ping"; onTriggered: vpnController.testLatency() }
    }

    Dialog {
        id: addServerDialog
        title: "Add Server"
        modal: true
        anchors.centerIn: parent
        background: Rectangle { color: cardColor; radius: 8 }
        ColumnLayout {
            width: parent.width
            spacing: 10
            TextField {
                id: serverUriInput
                Layout.fillWidth: true
                placeholderText: "vless://... or vmess://..."
                color: textColor
                background: Rectangle { color: bgColor; radius: 4 }
            }
            Button {
                text: "Add"
                Layout.alignment: Qt.AlignRight
                background: Rectangle { color: accentColor; radius: 4 }
                contentItem: Text { text: parent.text; color: "#FFFFFF" }
                onClicked: {
                    if (serverUriInput.text.length > 0) {
                        vpnController.addSubscription(serverUriInput.text)
                        addServerDialog.close()
                    }
                }
            }
        }
    }

    ColumnLayout {
        anchors.fill: parent
        anchors.margins: 20
        spacing: 15

        // Header
        RowLayout {
            Layout.fillWidth: true
            Text {
                text: "KINGO VPN"
                color: textColor
                font.pixelSize: 22
                font.bold: true
                font.letterSpacing: 1
            }
            Item { Layout.fillWidth: true }
            Button {
                text: "⚙"
                background: Rectangle { color: "transparent" }
                contentItem: Text { text: parent.text; color: subTextColor; font.pointSize: 20 }
                onClicked: settingsView.visible = true
            }
        }

        // Connection Status Card
        Rectangle {
            Layout.fillWidth: true
            height: 200
            color: bgColor
            radius: 12

            ColumnLayout {
                anchors.fill: parent
                anchors.margins: 15
                spacing: 5

                Text {
                    text: getStatusText()
                    color: getButtonColor()
                    font.pixelSize: 18
                    font.bold: true
                    Layout.alignment: Qt.AlignHCenter
                }
                Text {
                    text: getSubStatusText()
                    color: subTextColor
                    font.pixelSize: 12
                    Layout.alignment: Qt.AlignHCenter
                    visible: connStatus != "connected"
                }

                RowLayout {
                    Layout.fillWidth: true
                    Layout.topMargin: 5
                    spacing: 10
                    visible: connStatus == "connected"

                    ColumnLayout {
                        Layout.fillWidth: true
                        spacing: 2
                        Text { text: vpnController.connectionTime; color: textColor; font.pixelSize: 14; font.bold: true; Layout.alignment: Qt.AlignHCenter }
                        Text { text: "Duration"; color: subTextColor; font.pixelSize: 9; Layout.alignment: Qt.AlignHCenter }
                    }
                    ColumnLayout {
                        Layout.fillWidth: true
                        spacing: 2
                        Text { id: downSpeedText; text: "0.0 B/s"; color: successColor; font.pixelSize: 14; font.bold: true; Layout.alignment: Qt.AlignHCenter }
                        Text { text: "Down"; color: subTextColor; font.pixelSize: 9; Layout.alignment: Qt.AlignHCenter }
                    }
                    ColumnLayout {
                        Layout.fillWidth: true
                        spacing: 2
                        Text { id: upSpeedText; text: "0.0 B/s"; color: "#FFC107"; font.pixelSize: 14; font.bold: true; Layout.alignment: Qt.AlignHCenter }
                        Text { text: "Up"; color: subTextColor; font.pixelSize: 9; Layout.alignment: Qt.AlignHCenter }
                    }
                    ColumnLayout {
                        Layout.fillWidth: true
                        spacing: 2
                        Text { text: vpnController.ipAddress; color: textColor; font.pixelSize: 14; font.bold: true; Layout.alignment: Qt.AlignHCenter }
                        Text { text: "IP"; color: subTextColor; font.pixelSize: 9; Layout.alignment: Qt.AlignHCenter }
                    }
                }

                Rectangle {
                    Layout.alignment: Qt.AlignHCenter
                    Layout.topMargin: 10
                    width: 80
                    height: 80
                    radius: 40
                    color: "transparent"
                    border.width: 3
                    border.color: getButtonColor()
                    Behavior on border.color { ColorAnimation { duration: 300 } }

                    Rectangle {
                        anchors.centerIn: parent
                        width: 70
                        height: 70
                        radius: 35
                        color: getButtonColor()
                        Behavior on color { ColorAnimation { duration: 300 } }
                        scale: powerMouseArea.pressed ? 0.95 : 1.0
                        Behavior on scale { NumberAnimation { duration: 100 } }

                        Text {
                            anchors.centerIn: parent
                            text: "⏻" // Power Symbol
                            color: "#FFFFFF"
                            font.pixelSize: 30
                        }

                        MouseArea {
                            id: powerMouseArea
                            anchors.fill: parent
                            cursorShape: Qt.PointingHandCursor
                            onClicked: {
                                if (connStatus == "connected") {
                                    vpnController.disconnectVpn()
                                } else if (connStatus != "connecting") {
                                    prevDownload = 0
                                    prevUpload = 0
                                    prevTime = Date.now()
                                    vpnController.autoConnect()
                                }
                            }
                        }
                    }
                }
            }
        }

        // Server List Section
        RowLayout {
            Layout.fillWidth: true
            Text {
                text: "Server List"
                color: textColor
                font.pixelSize: 16
                font.bold: true
            }
            Item { Layout.fillWidth: true }
            Button {
                text: "⋮"
                background: Rectangle { color: "transparent" }
                contentItem: Text { text: parent.text; color: subTextColor; font.pointSize: 20 }
                onClicked: serverMenu.popup()
            }
        }

        // Tabs
        RowLayout {
            Layout.fillWidth: true
            spacing: 0
            Repeater {
                model: ["Favorites", "Servers", "Custom"]
                delegate: Rectangle {
                    Layout.fillWidth: true
                    height: 40
                    color: currentTab === index ? "#2C2C2C" : "transparent"
                    radius: 4
                    scale: tabMouseArea.pressed ? 0.95 : 1.0
                    Behavior on scale { NumberAnimation { duration: 100 } }
                    Text {
                        anchors.centerIn: parent
                        text: modelData
                        color: currentTab === index ? accentColor : subTextColor
                        font.bold: currentTab === index
                        font.pixelSize: 12
                    }
                    MouseArea {
                        id: tabMouseArea
                        anchors.fill: parent
                        cursorShape: Qt.PointingHandCursor
                        onClicked: currentTab = index
                    }
                }
            }
        }

        // Action Buttons
        RowLayout {
            Layout.fillWidth: true
            spacing: 10
            Button {
                Layout.fillWidth: true
                text: "Get Active Servers"
                background: Rectangle { color: "#2C2C2C"; radius: 6; border.color: "#333333"; border.width: 1 }
                contentItem: Text { text: parent.text; color: textColor; font.pixelSize: 10; font.bold: true; horizontalAlignment: Text.AlignHCenter; verticalAlignment: Text.AlignVCenter }
                onClicked: vpnController.addSubscription("https://raw.githubusercontent.com/MhdiTaheri/VpnHub/main/sub")
            }
            Button {
                Layout.fillWidth: true
                text: "Best Server"
                background: Rectangle { color: "#2C2C2C"; radius: 6; border.color: "#333333"; border.width: 1 }
                contentItem: Text { text: parent.text; color: textColor; font.pixelSize: 10; font.bold: true; horizontalAlignment: Text.AlignHCenter; verticalAlignment: Text.AlignVCenter }
                onClicked: {
                    prevDownload = 0; prevUpload = 0; prevTime = Date.now()
                    selectedServerID = "best"
                    vpnController.autoConnect()
                }
            }
        }

        // Server ListView
        ListView {
            id: serverListView
            Layout.fillWidth: true
            Layout.fillHeight: true
            clip: true
            model: currentTab === 0 ? vpnController.favorites : (currentTab === 1 ? vpnController.servers : vpnController.custom)
            spacing: 8
            ScrollBar.vertical: ScrollBar { active: true; policy: ScrollBar.AsNeeded }

            delegate: Rectangle {
                width: serverListView.width
                height: 55
                color: selectedServerID === modelData.id ? "#2C2C2C" : "#1E1E1E"
                radius: 6
                border.color: selectedServerID === modelData.id ? accentColor : "transparent"
                border.width: selectedServerID === modelData.id ? 1 : 0
                scale: delegateMouseArea.pressed ? 0.98 : 1.0
                Behavior on scale { NumberAnimation { duration: 100 } }

                RowLayout {
                    anchors.fill: parent
                    anchors.margins: 12
                    spacing: 10
                    Rectangle {
                        width: 8; height: 8; radius: 4
                        color: getPingColor(modelData.latency)
                        Layout.alignment: Qt.AlignVCenter
                    }
                    ColumnLayout {
                        Layout.fillWidth: true
                        spacing: 1
                        Text { text: modelData.name; color: "#FFFFFF"; font.bold: true; font.pixelSize: 13; elide: Text.ElideRight; Layout.fillWidth: true }
                        Text { text: modelData.protocol.toUpperCase() + " • " + modelData.address; color: "#B0B0B0"; font.pixelSize: 10 }
                    }
                    Text { 
                        text: modelData.latency == 9999 ? "N/A" : (modelData.latency == 0 ? "-" : modelData.latency + "ms")
                        color: getPingColor(modelData.latency)
                        font.bold: true
                        font.pixelSize: 11
                        Layout.alignment: Qt.AlignVCenter
                    }
                }

                MouseArea {
                    id: delegateMouseArea
                    anchors.fill: parent
                    cursorShape: Qt.PointingHandCursor
                    onClicked: {
                        selectedServerID = modelData.id
                        prevDownload = 0; prevUpload = 0; prevTime = Date.now()
                        vpnController.connectToServer(modelData.uri)
                    }
                }
            }
        }
    }

    // Settings Overlay
    Rectangle {
        id: settingsView
        anchors.fill: parent
        color: bgColor
        visible: false
        z: 100

        ColumnLayout {
            anchors.fill: parent
            anchors.margins: 20
            spacing: 15

            RowLayout {
                Layout.fillWidth: true
                Button {
                    text: "← Back"
                    background: Rectangle { color: "transparent" }
                    contentItem: Text { text: parent.text; color: accentColor; font.pixelSize: 14 }
                    onClicked: settingsView.visible = false
                }
                Item { Layout.fillWidth: true }
            }

            Text {
                text: "Split Tunneling"
                color: "#FFFFFF"
                font.pixelSize: 18
                font.bold: true
                Layout.topMargin: 10
            }
            RowLayout {
                Layout.fillWidth: true
                ColumnLayout {
                    Layout.fillWidth: true
                    Text { text: "Choose which apps use the VPN"; color: "#B0B0B0"; font.pixelSize: 12 }
                }
                Switch {
                    checked: false
                    onToggled: {
                        // Placeholder for Split Tunneling logic
                    }
                }
            }

            Rectangle { Layout.fillWidth: true; height: 1; color: "#333333"; Layout.topMargin: 10 }

            Text {
                text: "About"
                color: "#FFFFFF"
                font.pixelSize: 18
                font.bold: true
                Layout.topMargin: 10
            }
            Text { text: "Kingo VPN v0.5"; color: "#B0B0B0"; font.pixelSize: 12 }
            Text { text: "GitHub (Project source code)"; color: accentColor; font.pixelSize: 12; Layout.topMargin: 5 }
            Text { text: "Telegram (Join our channel)"; color: accentColor; font.pixelSize: 12; Layout.topMargin: 5 }
            
            Item { Layout.fillHeight: true }
            Text { text: "by Kingo team"; color: "#666666"; font.pixelSize: 10; Layout.alignment: Qt.AlignHCenter }
        }
    }
}
