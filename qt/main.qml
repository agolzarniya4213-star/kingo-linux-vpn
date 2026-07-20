import QtQuick
import QtQuick.Controls
import QtQuick.Layouts

ApplicationWindow {
    id: root
    width: 900
    height: 600
    minimumWidth: 800
    minimumHeight: 500
    visible: true
    title: "Kingo VPN"
    color: "#070D18"

    property string connStatus: vpnController ? vpnController.status : "disconnected"
    property int currentTab: 1
    property string selectedServerID: ""
    property real prevDownload: 0
    property real prevUpload: 0
    property real prevTime: 0
    property bool autoScanActive: false
    property int autoScanProgress: 0

    readonly property color cardColor: "#1A2636"
    readonly property color surfaceColor: "#0E1A2C"
    readonly property color borderColor: "#243245"
    readonly property color accentColor: "#18CFFF"
    readonly property color textColor: "#FFFFFF"
    readonly property color subTextColor: "#8899AA"
    readonly property color successColor: "#4CAF50"
    readonly property color warningColor: "#FFC107"
    readonly property color errorColor: "#F44336"

    Timer {
        id: trafficTimer
        interval: 1000
        running: root.connStatus == "connected"
        repeat: true
        onTriggered: vpnController.getTraffic()
    }

    Timer {
        id: scanTimer
        interval: 100
        repeat: true
        running: root.autoScanActive
        onTriggered: {
            if (root.autoScanProgress < 95) {
                root.autoScanProgress += 1
            }
        }
    }

    function formatSpeed(bytes) {
        if (!bytes || bytes < 0) return "0.0 B/s"
        if (bytes < 1024) return bytes.toFixed(1) + " B/s"
        if (bytes < 1048576) return (bytes / 1024).toFixed(1) + " KB/s"
        return (bytes / 1048576).toFixed(1) + " MB/s"
    }

    function getPingColor(latency) {
        if (latency == 9999 || latency == 0) return root.subTextColor
        if (latency < 150) return root.successColor
        if (latency < 300) return root.warningColor
        return root.errorColor
    }

    function getStatusColor() {
        if (root.connStatus == "connected") return root.successColor
        if (root.connStatus == "connecting") return root.warningColor
        return root.accentColor
    }
    
    function getStatusText() {
        if (root.connStatus == "connecting") return "Connecting..."
        if (root.connStatus == "connected") return "Connected"
        return "Disconnected"
    }

    onClosing: (close) => {
        close.accepted = false
        hide()
        trayIcon.showMessage("Kingo VPN", "Minimized to tray.")
    }

    Connections {
        target: trayIcon
        function onActivateRequested() { root.show(); root.raise(); root.requestActivate() }
        function onConnectRequested() { root.show(); root.raise() }
        function onDisconnectRequested() { vpnController.disconnectVpn() }
        function onQuitRequested() { Qt.quit() }
    }

    Connections {
        target: vpnController
        function onStatusChanged() {
            if (root.connStatus === "connected" || root.connStatus === "disconnected") {
                root.autoScanActive = false
            }
        }
        function onTrafficChanged() {
            var currTime = Date.now()
            var dt = (currTime - root.prevTime) / 1000.0
            if (dt > 0) {
                txtDownSpeed.text = formatSpeed((vpnController.downloadSpeed - root.prevDownload) / dt)
                txtUpSpeed.text = formatSpeed((vpnController.uploadSpeed - root.prevUpload) / dt)
            }
            root.prevDownload = vpnController.downloadSpeed
            root.prevUpload = vpnController.uploadSpeed
            root.prevTime = currTime
        }
    }

    Menu {
        id: serverMenu
        background: Rectangle {
            color: root.cardColor
            border.color: root.borderColor
            radius: 12
        }
        MenuItem { text: "Add Server"; onTriggered: addServerDialog.open() }
        MenuItem { text: "Fetch Servers"; onTriggered: vpnController.addSubscription("https://raw.githubusercontent.com/MhdiTaheri/VpnHub/main/sub") }
        MenuItem { text: "Test All Ping"; onTriggered: vpnController.testLatency() }
    }

    Dialog {
        id: addServerDialog
        title: "Add Server"
        modal: true
        anchors.centerIn: parent
        background: Rectangle { color: root.cardColor; radius: 22 }
        ColumnLayout {
            width: parent.width
            spacing: 12
            TextField {
                id: serverUriInput
                Layout.fillWidth: true
                placeholderText: "vless://... or vmess://..."
                color: root.textColor
                background: Rectangle { color: root.surfaceColor; radius: 12 }
            }
            Button {
                text: "Add"
                Layout.alignment: Qt.AlignRight
                background: Rectangle { color: root.accentColor; radius: 12 }
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

    // Desktop Split Layout
    RowLayout {
        anchors.fill: parent
        anchors.margins: 20
        spacing: 20

        // Left Column (Status & Stats)
        ColumnLayout {
            Layout.preferredWidth: 320
            Layout.fillHeight: true
            spacing: 16

            Text {
                text: "Kingo VPN"
                color: root.textColor
                font.pixelSize: 28
                font.weight: Font.Bold
                Layout.alignment: Qt.AlignHCenter
            }

            Text {
                text: getStatusText()
                color: getStatusColor()
                font.pixelSize: 18
                font.weight: Font.Bold
                Layout.alignment: Qt.AlignHCenter
            }
            
            Text {
                text: root.connStatus == "connected" ? "" : "Tap to connect"
                color: root.subTextColor
                font.pixelSize: 11
                opacity: 0.8
                Layout.alignment: Qt.AlignHCenter
            }

            // Power Button
            Item {
                Layout.alignment: Qt.AlignHCenter
                width: 160
                height: 160

                Rectangle {
                    anchors.fill: parent
                    radius: width / 2
                    color: "transparent"
                    border.color: getStatusColor()
                    border.width: 4
                    Behavior on border.color { ColorAnimation { duration: 300 } }
                }

                Rectangle {
                    id: innerCircle
                    anchors.centerIn: parent
                    width: 140
                    height: 140
                    radius: width / 2
                    color: root.surfaceColor
                    border.color: root.borderColor
                    border.width: 1
                    scale: powerMouseArea.pressed ? 0.95 : 1.0
                    Behavior on scale { NumberAnimation { duration: 100; easing.type: Easing.OutCubic } }

                    Rectangle {
                        anchors.fill: parent
                        radius: parent.radius
                        color: getStatusColor()
                        opacity: root.connStatus == "connecting" ? 0.3 : 0.0
                        z: -1
                        SequentialAnimation on scale {
                            running: root.connStatus == "connecting"
                            loops: Animation.Infinite
                            NumberAnimation { to: 1.2; duration: 800; easing.type: Easing.OutQuad }
                            NumberAnimation { to: 1.0; duration: 800; easing.type: Easing.InQuad }
                        }
                        SequentialAnimation on opacity {
                            running: root.connStatus == "connecting"
                            loops: Animation.Infinite
                            NumberAnimation { to: 0.0; duration: 800 }
                            NumberAnimation { to: 0.3; duration: 800 }
                        }
                    }

                    Text {
                        anchors.centerIn: parent
                        text: "⏻"
                        color: "#FFFFFF"
                        font.pixelSize: 50
                        visible: !root.autoScanActive
                    }
                    
                    ColumnLayout {
                        anchors.centerIn: parent
                        visible: root.autoScanActive
                        spacing: 2
                        Text { text: "Scanning..."; color: "#FFFFFF"; font.pixelSize: 12; font.weight: Font.Bold; Layout.alignment: Qt.AlignHCenter }
                        Text { text: root.autoScanProgress + "%"; color: root.accentColor; font.pixelSize: 16; font.weight: Font.Bold; Layout.alignment: Qt.AlignHCenter }
                    }

                    MouseArea {
                        id: powerMouseArea
                        anchors.fill: parent
                        cursorShape: Qt.PointingHandCursor
                        // FIX: Allow disconnecting while connecting
                        onClicked: {
                            if (root.connStatus == "connected" || root.connStatus == "connecting") {
                                root.autoScanActive = false
                                vpnController.disconnectVpn()
                            } else {
                                root.prevDownload = 0
                                root.prevUpload = 0
                                root.prevTime = Date.now()
                                root.autoScanActive = true
                                root.autoScanProgress = 0
                                vpnController.autoConnect()
                            }
                        }
                    }
                }
            }

            Text {
                text: "Best Server (Auto Connect)"
                color: root.accentColor
                font.pixelSize: 12
                font.weight: Font.Bold
                Layout.alignment: Qt.AlignHCenter
                MouseArea {
                    anchors.fill: parent
                    cursorShape: Qt.PointingHandCursor
                    onClicked: {
                        root.prevDownload = 0; root.prevUpload = 0; root.prevTime = Date.now()
                        root.selectedServerID = "best"
                        root.autoScanActive = true
                        root.autoScanProgress = 0
                        vpnController.autoConnect()
                    }
                }
            }

            // Stats Grid
            GridLayout {
                Layout.fillWidth: true
                columns: 2
                rowSpacing: 10
                columnSpacing: 10

                Rectangle {
                    Layout.fillWidth: true; height: 50; radius: 12; color: root.cardColor; border.color: root.borderColor
                    ColumnLayout { anchors.fill: parent; anchors.margins: 8; spacing: 1
                        Text { text: "Time"; color: root.subTextColor; font.pixelSize: 9 }
                        Text { text: vpnController ? vpnController.connectionTime : "00:00:00"; color: root.textColor; font.pixelSize: 14; font.weight: Font.Bold }
                    }
                }
                Rectangle {
                    Layout.fillWidth: true; height: 50; radius: 12; color: root.cardColor; border.color: root.borderColor
                    ColumnLayout { anchors.fill: parent; anchors.margins: 8; spacing: 1
                        Text { text: "VPN IP"; color: root.subTextColor; font.pixelSize: 9 }
                        Text { text: vpnController ? vpnController.ipAddress : "0.0.0.0"; color: root.textColor; font.pixelSize: 14; font.weight: Font.Bold; elide: Text.ElideRight }
                    }
                }
                Rectangle {
                    Layout.fillWidth: true; height: 50; radius: 12; color: root.cardColor; border.color: root.borderColor
                    ColumnLayout { anchors.fill: parent; anchors.margins: 8; spacing: 1
                        Text { text: "Download"; color: root.subTextColor; font.pixelSize: 9 }
                        Text { id: txtDownSpeed; text: "0.0 B/s"; color: root.successColor; font.pixelSize: 14; font.weight: Font.Bold }
                    }
                }
                Rectangle {
                    Layout.fillWidth: true; height: 50; radius: 12; color: root.cardColor; border.color: root.borderColor
                    ColumnLayout { anchors.fill: parent; anchors.margins: 8; spacing: 1
                        Text { text: "Upload"; color: root.subTextColor; font.pixelSize: 9 }
                        Text { id: txtUpSpeed; text: "0.0 B/s"; color: root.warningColor; font.pixelSize: 14; font.weight: Font.Bold }
                    }
                }
            }
            Item { Layout.fillHeight: true } // Spacer
        }

        // Right Column (Server List & Logs)
        ColumnLayout {
            Layout.fillWidth: true
            Layout.fillHeight: true
            spacing: 12

            // Server List Card
            Rectangle {
                Layout.fillWidth: true
                Layout.fillHeight: true
                color: root.cardColor
                radius: 18
                border.color: root.borderColor
                border.width: 1

                ColumnLayout {
                    anchors.fill: parent
                    anchors.margins: 16
                    spacing: 12

                    RowLayout {
                        Layout.fillWidth: true
                        Text { text: "Server List"; color: root.textColor; font.pixelSize: 20; font.weight: Font.Bold }
                        Item { Layout.fillWidth: true }
                        
                        // FIX: Menu Button using MouseArea
                        Item {
                            width: 30; height: 30
                            Text { anchors.centerIn: parent; text: "⋮"; color: root.subTextColor; font.pixelSize: 20 }
                            MouseArea { anchors.fill: parent; cursorShape: Qt.PointingHandCursor; onClicked: serverMenu.popup() }
                        }
                    }

                    // Tabs
                    Rectangle {
                        Layout.fillWidth: true; height: 36; radius: 12; color: root.surfaceColor; border.color: root.borderColor
                        RowLayout {
                            anchors.fill: parent; anchors.margins: 4; spacing: 4
                            Repeater {
                                model: ["Favorites", "Servers", "Custom"]
                                delegate: Rectangle {
                                    Layout.fillWidth: true; Layout.fillHeight: true; radius: 8
                                    color: root.currentTab === index ? root.cardColor : "transparent"
                                    border.color: root.currentTab === index ? root.accentColor : "transparent"
                                    border.width: 1
                                    scale: tabMouseArea.pressed ? 0.95 : 1.0
                                    Behavior on scale { NumberAnimation { duration: 100 } }
                                    Behavior on color { ColorAnimation { duration: 200 } }
                                    Text { anchors.centerIn: parent; text: modelData; color: root.currentTab === index ? root.textColor : root.subTextColor; font.weight: root.currentTab === index ? Font.Bold : Font.Normal; font.pixelSize: 11 }
                                    MouseArea { id: tabMouseArea; anchors.fill: parent; cursorShape: Qt.PointingHandCursor; onClicked: root.currentTab = index }
                                }
                            }
                        }
                    }

                    Rectangle {
                        Layout.fillWidth: true; height: 36; radius: 12
                        gradient: Gradient { orientation: Gradient.Horizontal; GradientStop { position: 0.0; color: root.accentColor }; GradientStop { position: 1.0; color: "#0FA8D9" } }
                        scale: activeBtnMouseArea.pressed ? 0.98 : 1.0
                        Behavior on scale { NumberAnimation { duration: 100 } }
                        Text { anchors.centerIn: parent; text: "Get Active Servers"; color: "#FFFFFF"; font.pixelSize: 12; font.weight: Font.Bold }
                        MouseArea { id: activeBtnMouseArea; anchors.fill: parent; cursorShape: Qt.PointingHandCursor; onClicked: { vpnController.addSubscription("https://raw.githubusercontent.com/MhdiTaheri/VpnHub/main/sub"); vpnController.fetchServers() } }
                    }

                    ListView {
                        id: serverListView
                        Layout.fillWidth: true; Layout.fillHeight: true; clip: true
                        model: root.currentTab === 0 ? (vpnController ? vpnController.favorites : []) : (root.currentTab === 1 ? (vpnController ? vpnController.servers : []) : (vpnController ? vpnController.custom : []))
                        spacing: 8
                        ScrollBar.vertical: ScrollBar { policy: ScrollBar.AsNeeded }

                        delegate: Rectangle {
                            width: serverListView.width; height: 50; radius: 12
                            color: root.selectedServerID === modelData.id ? "#243245" : root.surfaceColor
                            border.color: root.selectedServerID === modelData.id ? root.accentColor : "transparent"
                            border.width: 1
                            scale: delegateMouseArea.pressed ? 0.98 : 1.0
                            Behavior on scale { NumberAnimation { duration: 100 } }
                            Behavior on color { ColorAnimation { duration: 200 } }

                            RowLayout {
                                anchors.fill: parent; anchors.margins: 12; spacing: 10
                                Rectangle { width: 8; height: 8; radius: 4; color: getPingColor(modelData.latency); Layout.alignment: Qt.AlignVCenter }
                                ColumnLayout {
                                    Layout.fillWidth: true; spacing: 1
                                    Text { text: modelData.name; color: root.textColor; font.weight: Font.Bold; font.pixelSize: 13; elide: Text.ElideRight; Layout.fillWidth: true }
                                    Text { text: modelData.protocol.toUpperCase() + " • " + modelData.address; color: root.subTextColor; font.pixelSize: 9 }
                                }
                                Text { text: modelData.latency == 9999 ? "N/A" : (modelData.latency == 0 ? "-" : modelData.latency + "ms"); color: getPingColor(modelData.latency); font.bold: true; font.pixelSize: 11; Layout.alignment: Qt.AlignVCenter }
                            }

                            MouseArea {
                                id: delegateMouseArea
                                anchors.fill: parent
                                cursorShape: Qt.PointingHandCursor
                                onClicked: {
                                    root.selectedServerID = modelData.id
                                    root.prevDownload = 0; root.prevUpload = 0; root.prevTime = Date.now()
                                    root.autoScanActive = false
                                    vpnController.connectToServer(modelData.uri)
                                }
                            }
                        }
                    }
                }
            }

            // Logs Panel
            Rectangle {
                Layout.fillWidth: true
                Layout.preferredHeight: 120
                color: root.surfaceColor
                radius: 12
                border.color: root.borderColor
                border.width: 1

                ColumnLayout {
                    anchors.fill: parent
                    anchors.margins: 8
                    spacing: 4

                    RowLayout {
                        Layout.fillWidth: true
                        Text { text: "Logs"; color: root.subTextColor; font.pixelSize: 10; font.bold: true }
                        Item { Layout.fillWidth: true }
                        Button {
                            text: "Copy"
                            implicitHeight: 20
                            background: Rectangle { color: root.cardColor; radius: 4; border.color: root.borderColor; border.width: 1 }
                            contentItem: Text { text: parent.text; color: root.accentColor; font.pixelSize: 9; font.bold: true; horizontalAlignment: Text.AlignHCenter; verticalAlignment: Text.AlignVCenter }
                            onClicked: vpnController.copyLogs()
                        }
                    }

                    ScrollView {
                        Layout.fillWidth: true
                        Layout.fillHeight: true
                        clip: true
                        TextArea {
                            text: vpnController.logs
                            color: root.subTextColor
                            font.family: "Monospace"
                            font.pixelSize: 9
                            wrapMode: TextArea.WrapAnywhere
                            readOnly: true
                            selectByMouse: true
                            background: Rectangle { color: "transparent" }
                        }
                    }
                }
            }
        }
    }
}
