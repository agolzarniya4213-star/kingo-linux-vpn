import QtQuick
import QtQuick.Controls
import QtQuick.Layouts

Window {
    id: mainWindow
    width: 420
    height: 850
    visible: true
    title: "Kingo VPN"
    color: "#0D1117" // GitHub Dark Background

    onClosing: (close) => {
        close.accepted = false
        hide()
        trayIcon.showMessage("Kingo VPN", "Minimized to tray.")
    }

    property string connStatus: vpnController ? vpnController.status : "loading"
    property var serverList: vpnController ? vpnController.servers : []
    property real prevDownload: 0
    property real prevUpload: 0
    property real prevTime: 0

    // GitHub Dark Theme Palette
    readonly property color bgColor: "#0D1117"
    readonly property color cardColor: "#161B22"
    readonly property color borderColor: "#30363D"
    readonly property color accentColor: "#58A6FF" // GitHub Blue
    readonly property color textColor: "#C9D1D9"
    readonly property color subTextColor: "#8B949E"
    readonly property color successColor: "#3FB950" // GitHub Green
    readonly property color errorColor: "#F85149"   // GitHub Red
    readonly property color warningColor: "#D29922" // GitHub Yellow

    Timer {
        interval: 1000
        running: connStatus == "connected"
        repeat: true
        onTriggered: vpnController.getTraffic()
    }

    function formatSpeed(bytes) {
        if (!bytes || bytes < 0) return "0 B/s"
        if (bytes < 1024) return bytes.toFixed(0) + " B/s"
        if (bytes < 1048576) return (bytes / 1024).toFixed(1) + " KB/s"
        return (bytes / 1048576).toFixed(2) + " MB/s"
    }

    function getPingColor(latency) {
        if (latency == 9999 || latency == 0) return subTextColor
        if (latency < 150) return successColor
        if (latency < 300) return warningColor
        return errorColor
    }

    function getButtonColor() {
        if (connStatus == "connected") return errorColor
        if (connStatus == "connecting") return warningColor
        return accentColor
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
                downloadSpeedText.text = formatSpeed((vpnController.downloadSpeed - prevDownload) / dt)
                uploadSpeedText.text = formatSpeed((vpnController.uploadSpeed - prevUpload) / dt)
            }
            prevDownload = vpnController.downloadSpeed
            prevUpload = vpnController.uploadSpeed
            prevTime = currTime
        }
    }

    // Reusable Button Component
    component ActionButton : Button {
        id: btn
        scale: pressed ? 0.95 : 1.0
        Behavior on scale { NumberAnimation { duration: 100; easing.type: Easing.OutCubic } }
        contentItem: Text {
            text: btn.text
            color: btn.textColor
            font.pixelSize: 11
            font.bold: true
            horizontalAlignment: Text.AlignHCenter
            verticalAlignment: Text.AlignVCenter
        }
    }

    ColumnLayout {
        anchors.fill: parent
        anchors.margins: 20
        spacing: 12

        // Header
        RowLayout {
            Layout.fillWidth: true
            Text {
                text: "KINGO VPN"
                color: textColor
                font.pixelSize: 20
                font.bold: true
                font.letterSpacing: 2
            }
            Item { Layout.fillWidth: true }
            Text {
                text: "v1.8"
                color: subTextColor
                font.pixelSize: 12
            }
        }

        // Circular Connect Button with Glows
        Rectangle {
            Layout.alignment: Qt.AlignHCenter
            Layout.topMargin: 10
            width: 160
            height: 160
            radius: 80
            color: "transparent"
            border.width: 4
            border.color: getButtonColor()
            Behavior on border.color { ColorAnimation { duration: 300 } }

            Rectangle {
                id: innerCircle
                anchors.centerIn: parent
                width: 140
                height: 140
                radius: 70
                color: getButtonColor()
                Behavior on color { ColorAnimation { duration: 300 } }
                scale: connectMouseArea.pressed ? 0.95 : 1.0
                Behavior on scale { NumberAnimation { duration: 100; easing.type: Easing.OutCubic } }

                Rectangle {
                    anchors.fill: parent
                    radius: parent.radius
                    color: parent.color
                    opacity: connStatus == "connecting" ? 0.4 : 0.0
                    z: -1
                    SequentialAnimation on scale {
                        running: connStatus == "connecting"
                        loops: Animation.Infinite
                        NumberAnimation { to: 1.3; duration: 600; easing.type: Easing.OutQuad }
                        NumberAnimation { to: 1.0; duration: 600; easing.type: Easing.InQuad }
                    }
                    SequentialAnimation on opacity {
                        running: connStatus == "connecting"
                        loops: Animation.Infinite
                        NumberAnimation { to: 0.0; duration: 600 }
                        NumberAnimation { to: 0.4; duration: 600 }
                    }
                }

                Text {
                    anchors.centerIn: parent
                    text: connStatus == "connecting" ? "..." : (connStatus == "connected" ? "DISCONNECT" : "CONNECT")
                    color: "#FFFFFF"
                    font.pixelSize: 16
                    font.bold: true
                }

                MouseArea {
                    id: connectMouseArea
                    anchors.fill: parent
                    cursorShape: Qt.PointingHandCursor
                    enabled: connStatus != "connecting"
                    onClicked: {
                        if (connStatus == "connected") {
                            vpnController.disconnectVpn()
                        } else {
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
            Layout.fillWidth: true
            text: connStatus.toUpperCase()
            color: textColor
            font.pixelSize: 16
            font.bold: true
            font.letterSpacing: 1
            horizontalAlignment: Text.AlignHCenter
        }

        // Speed Stats
        RowLayout {
            Layout.fillWidth: true
            spacing: 12
            visible: connStatus == "connected"

            Rectangle {
                Layout.fillWidth: true
                height: 45
                color: cardColor
                radius: 8
                border.color: borderColor
                border.width: 1
                ColumnLayout {
                    anchors.centerIn: parent
                    Text { text: "DOWNLOAD"; color: subTextColor; font.pixelSize: 8; font.bold: true; Layout.alignment: Qt.AlignHCenter }
                    Text { id: downloadSpeedText; text: "0 B/s"; color: successColor; font.pixelSize: 14; font.bold: true; Layout.alignment: Qt.AlignHCenter }
                }
            }

            Rectangle {
                Layout.fillWidth: true
                height: 45
                color: cardColor
                radius: 8
                border.color: borderColor
                border.width: 1
                ColumnLayout {
                    anchors.centerIn: parent
                    Text { text: "UPLOAD"; color: subTextColor; font.pixelSize: 8; font.bold: true; Layout.alignment: Qt.AlignHCenter }
                    Text { id: uploadSpeedText; text: "0 B/s"; color: textColor; font.pixelSize: 14; font.bold: true; Layout.alignment: Qt.AlignHCenter }
                }
            }
        }

        // Subscription Management
        RowLayout {
            Layout.fillWidth: true
            spacing: 8
            TextField {
                id: subUrlField
                Layout.fillWidth: true
                height: 40
                placeholderText: "Subscription URL..."
                color: textColor
                font.pixelSize: 12
                text: "https://raw.githubusercontent.com/MhdiTaheri/VpnHub/main/sub"
                background: Rectangle { color: cardColor; radius: 6; border.color: subUrlField.activeFocus ? accentColor : borderColor; border.width: 1 }
            }
            
            ActionButton {
                text: "Update"
                property color textColor: accentColor
                implicitWidth: 70
                implicitHeight: 40
                background: Rectangle { color: cardColor; radius: 6; border.color: accentColor; border.width: 1 }
                onClicked: { if (subUrlField.text.length > 0) vpnController.addSubscription(subUrlField.text) }
            }

            ActionButton {
                text: "Clear"
                property color textColor: errorColor
                implicitWidth: 50
                implicitHeight: 40
                background: Rectangle { color: cardColor; radius: 6; border.color: errorColor; border.width: 1 }
                onClicked: { vpnController.clearServers() }
            }
        }

        // Server List
        ListView {
            id: serverListView
            Layout.fillWidth: true
            Layout.fillHeight: true
            clip: true
            model: serverList
            spacing: 8
            ScrollBar.vertical: ScrollBar { active: true; policy: ScrollBar.AsNeeded }

            header: Rectangle {
                width: serverListView.width
                height: 55
                color: accentColor
                radius: 8
                scale: headerMouseArea.pressed ? 0.98 : 1.0
                Behavior on scale { NumberAnimation { duration: 100 } }
                RowLayout {
                    anchors.fill: parent
                    anchors.margins: 12
                    spacing: 10
                    Rectangle { width: 8; height: 8; radius: 4; color: "#FFFFFF" }
                    ColumnLayout {
                        Layout.fillWidth: true
                        spacing: 1
                        Text { text: "Best Server (Auto Connect)"; color: "#FFFFFF"; font.bold: true; font.pixelSize: 13 }
                        Text { text: "FIND THE FASTEST SERVER"; color: "#E2E8F0"; font.pixelSize: 8 }
                    }
                }
                MouseArea {
                    id: headerMouseArea
                    anchors.fill: parent
                    cursorShape: Qt.PointingHandCursor
                    onClicked: {
                        prevDownload = 0; prevUpload = 0; prevTime = Date.now()
                        vpnController.autoConnect()
                    }
                }
            }

            delegate: Rectangle {
                width: serverListView.width
                height: 50
                color: cardColor
                radius: 8
                border.color: borderColor
                border.width: 1
                scale: delegateMouseArea.pressed ? 0.98 : 1.0
                Behavior on scale { NumberAnimation { duration: 100 } }
                RowLayout {
                    anchors.fill: parent
                    anchors.margins: 12
                    spacing: 10
                    Rectangle { width: 8; height: 8; radius: 4; color: getPingColor(modelData.latency) }
                    ColumnLayout {
                        Layout.fillWidth: true
                        spacing: 1
                        Text { text: modelData.name; color: textColor; font.bold: true; font.pixelSize: 12; elide: Text.ElideRight; Layout.fillWidth: true }
                        Text { text: modelData.protocol.toUpperCase() + " • " + modelData.address; color: subTextColor; font.pixelSize: 9 }
                    }
                    Text { text: modelData.latency == 9999 ? "N/A" : (modelData.latency == 0 ? "-" : modelData.latency + "ms"); color: getPingColor(modelData.latency); font.bold: true; font.pixelSize: 11 }
                }
                MouseArea {
                    id: delegateMouseArea
                    anchors.fill: parent
                    cursorShape: Qt.PointingHandCursor
                    onClicked: {
                        prevDownload = 0; prevUpload = 0; prevTime = Date.now()
                        vpnController.connectToServer(modelData.uri)
                    }
                }
            }
        }

        // Logs Panel
        RowLayout {
            Layout.fillWidth: true
            spacing: 8
            Text { text: "LOGS"; color: subTextColor; font.pixelSize: 10; font.bold: true; font.letterSpacing: 1 }
            Item { Layout.fillWidth: true }
            
            ActionButton {
                text: "Copy Logs"
                property color textColor: accentColor
                implicitHeight: 25
                background: Rectangle { color: cardColor; radius: 4; border.color: borderColor; border.width: 1 }
                onClicked: { vpnController.copyLogs() }
            }
        }

        Rectangle {
            Layout.fillWidth: true
            height: 100
            color: "#010409" // Deep console black
            radius: 8
            border.color: borderColor
            border.width: 1
            clip: true

            ScrollView {
                anchors.fill: parent
                anchors.margins: 8
                ScrollBar.vertical.policy: ScrollBar.AsNeeded
                
                TextArea {
                    text: vpnController.logs
                    color: subTextColor
                    font.family: "Monospace"
                    font.pixelSize: 10
                    wrapMode: TextArea.WrapAnywhere
                    readOnly: true
                    selectByMouse: true
                    background: Rectangle { color: "transparent" }
                }
            }
        }
    }
}
