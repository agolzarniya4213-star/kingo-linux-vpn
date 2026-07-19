import QtQuick
import QtQuick.Controls
import QtQuick.Layouts

ApplicationWindow {
    id: mainWindow
    width: 480
    height: 860
    visible: true
    title: "Kingo VPN"
    color: "#070D18"

    onClosing: (close) => {
        close.accepted = false
        hide()
        trayIcon.showMessage("Kingo VPN", "Minimized to tray.")
    }

    Item {
        id: container
        anchors.centerIn: parent
        width: Math.min(780, parent.width)
        height: parent.height

        property string connStatus: vpnController ? vpnController.status : "loading"
        property int currentTab: 1
        property string selectedServerID: ""
        property real prevDownload: 0
        property real prevUpload: 0
        property real prevTime: 0

        // FIX: Removed 'readonly' to prevent qmlimportscanner syntax errors
        property color cardColor: "#1A2636"
        property color cardBorderColor: "#243245"
        property color accentColor: "#18CFFF"
        property color textColor: "#FFFFFF"
        property color subTextColor: "#8899AA"
        property color successColor: "#4CAF50"

        Timer {
            interval: 1000
            running: container.connStatus == "connected"
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
            if (container.connStatus == "connected") return successColor
            if (container.connStatus == "connecting") return "#FFC107"
            return accentColor
        }
        
        function getStatusText() {
            if (container.connStatus == "connecting") return "Connecting..."
            if (container.connStatus == "connected") return "Connected"
            return "Disconnected"
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
            background: Rectangle { color: cardColor; border.color: cardBorderColor; radius: 12 }
            MenuItem { text: "Add Server"; onTriggered: addServerDialog.open() }
            MenuItem { text: "Fetch Servers"; onTriggered: vpnController.addSubscription("https://raw.githubusercontent.com/MhdiTaheri/VpnHub/main/sub") }
            MenuItem { text: "Test All Ping"; onTriggered: vpnController.testLatency() }
        }

        Dialog {
            id: addServerDialog
            title: "Add Server"
            modal: true
            anchors.centerIn: parent
            background: Rectangle { color: cardColor; radius: 22 }
            ColumnLayout {
                width: parent.width
                spacing: 12
                TextField {
                    id: serverUriInput
                    Layout.fillWidth: true
                    placeholderText: "vless://... or vmess://..."
                    color: textColor
                    background: Rectangle { color: "#0E1A2C"; radius: 12 }
                }
                Button {
                    text: "Add"
                    Layout.alignment: Qt.AlignRight
                    background: Rectangle { color: accentColor; radius: 12 }
                    contentItem: Text { text: parent.text; color: "#FFFFFF"; font.bold: true }
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
            anchors.margins: 24
            spacing: 16

            Text {
                text: "Kingo VPN"
                color: textColor
                font.pixelSize: 38
                font.weight: Font.Bold
                Layout.alignment: Qt.AlignHCenter
            }

            Text {
                text: getStatusText()
                color: getButtonColor()
                font.pixelSize: 22
                font.weight: Font.Bold
                Layout.alignment: Qt.AlignHCenter
            }
            Text {
                text: container.connStatus == "connected" ? "" : "Tap to connect"
                color: subTextColor
                font.pixelSize: 12
                opacity: 0.8
                Layout.alignment: Qt.AlignHCenter
            }

            Item {
                Layout.alignment: Qt.AlignHCenter
                Layout.topMargin: 10
                width: 180
                height: 180

                Rectangle {
                    anchors.fill: parent
                    radius: width / 2
                    color: "transparent"
                    border.color: getButtonColor()
                    border.width: 4
                    Behavior on border.color { ColorAnimation { duration: 300 } }
                }

                Rectangle {
                    id: innerCircle
                    anchors.centerIn: parent
                    width: 160
                    height: 160
                    radius: width / 2
                    color: "#0E1A2C"
                    border.color: cardBorderColor
                    border.width: 1
                    scale: powerMouseArea.pressed ? 0.95 : 1.0
                    Behavior on scale { NumberAnimation { duration: 100; easing.type: Easing.OutCubic } }

                    Rectangle {
                        anchors.fill: parent
                        radius: parent.radius
                        color: getButtonColor()
                        opacity: container.connStatus == "connecting" ? 0.3 : 0.0
                        z: -1
                        SequentialAnimation on scale {
                            running: container.connStatus == "connecting"
                            loops: Animation.Infinite
                            NumberAnimation { to: 1.2; duration: 800; easing.type: Easing.OutQuad }
                            NumberAnimation { to: 1.0; duration: 800; easing.type: Easing.InQuad }
                        }
                        SequentialAnimation on opacity {
                            running: container.connStatus == "connecting"
                            loops: Animation.Infinite
                            NumberAnimation { to: 0.0; duration: 800 }
                            NumberAnimation { to: 0.3; duration: 800 }
                        }
                    }

                    Text {
                        anchors.centerIn: parent
                        text: "⏻"
                        color: "#FFFFFF"
                        font.pixelSize: 60
                    }

                    MouseArea {
                        id: powerMouseArea
                        anchors.fill: parent
                        cursorShape: Qt.PointingHandCursor
                        onClicked: {
                            if (container.connStatus == "connected") {
                                vpnController.disconnectVpn()
                            } else if (container.connStatus != "connecting") {
                                prevDownload = 0
                                prevUpload = 0
                                prevTime = Date.now()
                                vpnController.autoConnect()
                            }
                        }
                    }
                }
            }

            Text {
                text: "Best Server (Auto Connect)"
                color: accentColor
                font.pixelSize: 14
                font.weight: Font.Bold
                Layout.alignment: Qt.AlignHCenter
                Layout.topMargin: 10
                MouseArea {
                    anchors.fill: parent
                    cursorShape: Qt.PointingHandCursor
                    onClicked: {
                        prevDownload = 0; prevUpload = 0; prevTime = Date.now()
                        selectedServerID = "best"
                        vpnController.autoConnect()
                    }
                }
            }

            Rectangle {
                Layout.fillWidth: true
                Layout.fillHeight: true
                color: cardColor
                radius: 22
                border.color: cardBorderColor
                border.width: 1

                ColumnLayout {
                    anchors.fill: parent
                    anchors.margins: 24
                    spacing: 16

                    RowLayout {
                        Layout.fillWidth: true
                        Text {
                            text: "Server List"
                            color: textColor
                            font.pixelSize: 28
                            font.weight: Font.Bold
                        }
                        Item { Layout.fillWidth: true }
                        Button {
                            text: "⋮"
                            background: Rectangle { color: "transparent" }
                            contentItem: Text { text: parent.text; color: subTextColor; font.pixelSize: 24 }
                            onClicked: serverMenu.popup()
                        }
                    }

                    Rectangle {
                        Layout.fillWidth: true
                        height: 44
                        color: "#0E1A2C"
                        radius: 18

                        RowLayout {
                            anchors.fill: parent
                            anchors.margins: 4
                            spacing: 4

                            Repeater {
                                model: ["Favorites", "Servers", "Custom"]
                                delegate: Rectangle {
                                    Layout.fillWidth: true
                                    Layout.fillHeight: true
                                    radius: 14
                                    color: container.currentTab === index ? cardColor : "transparent"
                                    border.color: container.currentTab === index ? accentColor : "transparent"
                                    border.width: 1
                                    scale: tabMouseArea.pressed ? 0.95 : 1.0
                                    Behavior on scale { NumberAnimation { duration: 100 } }
                                    Behavior on color { ColorAnimation { duration: 200 } }

                                    Text {
                                        anchors.centerIn: parent
                                        text: modelData
                                        color: container.currentTab === index ? textColor : subTextColor
                                        font.weight: container.currentTab === index ? Font.Bold : Font.Normal
                                        font.pixelSize: 12
                                    }

                                    MouseArea {
                                        id: tabMouseArea
                                        anchors.fill: parent
                                        cursorShape: Qt.PointingHandCursor
                                        onClicked: container.currentTab = index
                                    }
                                }
                            }
                        }
                    }

                    Rectangle {
                        Layout.fillWidth: true
                        height: 48
                        radius: 18
                        gradient: Gradient {
                            orientation: Gradient.Horizontal
                            GradientStop { position: 0.0; color: "#18CFFF" }
                            GradientStop { position: 1.0; color: "#0FA8D9" }
                        }
                        scale: activeBtnMouseArea.pressed ? 0.98 : 1.0
                        Behavior on scale { NumberAnimation { duration: 100 } }

                        Text {
                            anchors.centerIn: parent
                            text: "Get Active Servers"
                            color: "#FFFFFF"
                            font.pixelSize: 14
                            font.weight: Font.Bold
                        }

                        MouseArea {
                            id: activeBtnMouseArea
                            anchors.fill: parent
                            cursorShape: Qt.PointingHandCursor
                            onClicked: vpnController.addSubscription("https://raw.githubusercontent.com/MhdiTaheri/VpnHub/main/sub")
                        }
                    }

                    ListView {
                        id: serverListView
                        Layout.fillWidth: true
                        Layout.fillHeight: true
                        clip: true
                        model: container.currentTab === 0 ? vpnController.favorites : (container.currentTab === 1 ? vpnController.servers : vpnController.custom)
                        spacing: 10
                        ScrollBar.vertical: ScrollBar { active: true; policy: ScrollBar.AsNeeded }

                        delegate: Rectangle {
                            width: serverListView.width
                            height: 60
                            radius: 16
                            color: selectedServerID === modelData.id ? "#243245" : "#0E1A2C"
                            border.color: selectedServerID === modelData.id ? accentColor : "transparent"
                            border.width: 1
                            scale: delegateMouseArea.pressed ? 0.98 : 1.0
                            Behavior on scale { NumberAnimation { duration: 100 } }
                            Behavior on color { ColorAnimation { duration: 200 } }

                            RowLayout {
                                anchors.fill: parent
                                anchors.margins: 16
                                spacing: 12

                                Rectangle {
                                    width: 10; height: 10; radius: 5
                                    color: getPingColor(modelData.latency)
                                    Layout.alignment: Qt.AlignVCenter
                                }

                                ColumnLayout {
                                    Layout.fillWidth: true
                                    spacing: 2
                                    Text { text: modelData.name; color: textColor; font.weight: Font.Bold; font.pixelSize: 14; elide: Text.ElideRight; Layout.fillWidth: true }
                                    Text { text: modelData.protocol.toUpperCase() + " • " + modelData.address; color: subTextColor; font.pixelSize: 10 }
                                }

                                Text { 
                                    text: modelData.latency == 9999 ? "N/A" : (modelData.latency == 0 ? "-" : modelData.latency + "ms")
                                    color: getPingColor(modelData.latency)
                                    font.bold: true
                                    font.pixelSize: 12
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
            }
        }

        Rectangle {
            id: settingsView
            anchors.fill: parent
            color: "#070D18"
            visible: false
            z: 100

            ColumnLayout {
                anchors.fill: parent
                anchors.margins: 24
                spacing: 20

                RowLayout {
                    Layout.fillWidth: true
                    Button {
                        text: "← Back"
                        background: Rectangle { color: "transparent" }
                        contentItem: Text { text: parent.text; color: accentColor; font.pixelSize: 16 }
                        onClicked: settingsView.visible = false
                    }
                    Item { Layout.fillWidth: true }
                }

                Text {
                    text: "Split Tunneling"
                    color: textColor
                    font.pixelSize: 24
                    font.weight: Font.Bold
                    Layout.topMargin: 20
                }
                RowLayout {
                    Layout.fillWidth: true
                    Text { text: "Choose which apps use the VPN"; color: subTextColor; font.pixelSize: 14; Layout.fillWidth: true }
                    Switch { checked: false; onToggled: { } }
                }

                Rectangle { Layout.fillWidth: true; height: 1; color: "#243245"; Layout.topMargin: 20 }

                Text {
                    text: "About"
                    color: textColor
                    font.pixelSize: 24
                    font.weight: Font.Bold
                    Layout.topMargin: 20
                }
                Text { text: "Kingo VPN v0.5"; color: subTextColor; font.pixelSize: 14 }
                Text { text: "GitHub (Project source code)"; color: accentColor; font.pixelSize: 14; Layout.topMargin: 10 }
                Text { text: "Telegram (Join our channel)"; color: accentColor; font.pixelSize: 14; Layout.topMargin: 5 }
                
                Item { Layout.fillHeight: true }
                Text { text: "by Kingo team"; color: "#556677"; font.pixelSize: 12; Layout.alignment: Qt.AlignHCenter }
            }
        }
    }
}
