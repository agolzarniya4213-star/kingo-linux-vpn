import QtQuick
import QtQuick.Controls
import QtQuick.Layouts

Window {
    id: mainWindow
    width: 420
    height: 780
    visible: true
    title: "Kingo VPN"
    color: "#0A0A0F" // Deep Cyber Black

    onClosing: (close) => {
        close.accepted = false
        hide()
        trayIcon.showMessage("Kingo VPN", "Application minimized to tray.")
    }

    // Defensive checks to prevent null crashes during initialization
    property string connStatus: vpnController ? vpnController.status : "loading"
    property var serverList: vpnController ? vpnController.servers : []
    property real prevDownload: 0
    property real prevUpload: 0
    property real prevTime: 0

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

    Connections {
        target: trayIcon
        function onActivateRequested() { mainWindow.show(); mainWindow.raise(); mainWindow.requestActivate() }
        function onConnectRequested() { mainWindow.show(); mainWindow.raise() }
        function onDisconnectRequested() { vpnController.disconnectVpn() }
        function onQuitRequested() { Qt.quit() }
    }

    Connections {
        target: vpnController
        function onErrorOccurred(error) {
            errorText.text = error
            errorBox.visible = true
            errorTimer.start()
        }
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

    ColumnLayout {
        anchors.fill: parent
        anchors.margins: 24
        spacing: 18

        // Header
        RowLayout {
            Layout.fillWidth: true
            Layout.topMargin: 10

            Rectangle {
                width: 32; height: 32
                color: "#00FF9D"
                radius: 6
                Text { text: "K"; color: "#0A0A0F"; font.bold: true; font.pixelSize: 20; anchors.centerIn: parent }
            }
            
            Text {
                text: "KINGO VPN"
                color: "#FFFFFF"
                font.pixelSize: 20
                font.bold: true
                font.letterSpacing: 2
                Layout.leftMargin: 10
            }
            
            Item { Layout.fillWidth: true }
            
            Text {
                text: "v2.0"
                color: "#4A4A5A"
                font.pixelSize: 12
            }
        }

        // Main Status Card
        Rectangle {
            Layout.fillWidth: true
            height: 220
            color: "#14141B"
            radius: 16
            border.color: "#22222E"
            border.width: 1

            ColumnLayout {
                anchors.fill: parent
                anchors.margins: 20
                spacing: 12

                Text {
                    text: "STATUS"
                    color: "#5A5A6A"
                    font.pixelSize: 11
                    font.bold: true
                    font.letterSpacing: 1
                }

                Text {
                    text: connStatus.toUpperCase()
                    color: connStatus == "connected" ? "#00FF9D" : (connStatus == "connecting" ? "#FFCC00" : "#FF3366")
                    font.pixelSize: 32
                    font.bold: true
                    font.letterSpacing: 1
                }

                // Speed Stats Row
                RowLayout {
                    Layout.fillWidth: true
                    Layout.topMargin: 10
                    spacing: 20
                    visible: connStatus == "connected"

                    Rectangle {
                        Layout.fillWidth: true
                        height: 60
                        color: "#1A1A24"
                        radius: 10
                        
                        ColumnLayout {
                            anchors.centerIn: parent
                            Text { text: "DOWNLOAD"; color: "#5A5A6A"; font.pixelSize: 9; font.bold: true; Layout.alignment: Qt.AlignHCenter }
                            Text { id: downloadSpeedText; text: "0 B/s"; color: "#00FF9D"; font.pixelSize: 16; font.bold: true; Layout.alignment: Qt.AlignHCenter }
                        }
                    }

                    Rectangle {
                        Layout.fillWidth: true
                        height: 60
                        color: "#1A1A24"
                        radius: 10
                        
                        ColumnLayout {
                            anchors.centerIn: parent
                            Text { text: "UPLOAD"; color: "#5A5A6A"; font.pixelSize: 9; font.bold: true; Layout.alignment: Qt.AlignHCenter }
                            Text { id: uploadSpeedText; text: "0 B/s"; color: "#FFFFFF"; font.pixelSize: 16; font.bold: true; Layout.alignment: Qt.AlignHCenter }
                        }
                    }
                }
            }
        }

        // Connect Button (Modern Hexagon-like Capsule)
        Button {
            Layout.fillWidth: true
            height: 56
            enabled: connStatus != "connecting"
            
            background: Rectangle {
                color: connStatus == "connected" ? "#FF3366" : "#00FF9D"
                radius: 12
            }
            contentItem: Text {
                text: connStatus == "connecting" ? "INITIALIZING..." : (connStatus == "connected" ? "DISCONNECT" : "CONNECT")
                color: "#0A0A0F"
                font.pixelSize: 16
                font.bold: true
                font.letterSpacing: 2
                horizontalAlignment: Text.AlignHCenter
                verticalAlignment: Text.AlignVCenter
            }
            onClicked: {
                if (connStatus == "connected") {
                    vpnController.disconnectVpn()
                } else {
                    prevDownload = 0; prevUpload = 0; prevTime = Date.now()
                    vpnController.autoConnect()
                }
            }
        }

        // Error Box
        Rectangle {
            id: errorBox
            Layout.fillWidth: true
            height: 40
            color: "#33000000"
            radius: 8
            visible: false
            border.color: "#FF3366"
            border.width: 1
            
            Text { id: errorText; anchors.centerIn: parent; color: "#FF3366"; font.bold: true; font.pixelSize: 12 }
            Timer { id: errorTimer; interval: 5000; onTriggered: errorBox.visible = false }
        }

        // Subscription Input
        TextField {
            id: subUrlField
            Layout.fillWidth: true
            height: 45
            placeholderText: "Enter Subscription URL..."
            color: "#FFFFFF"
            font.pixelSize: 13
            text: "https://raw.githubusercontent.com/MhdiTaheri/VpnHub/main/sub"
            background: Rectangle {
                color: "#14141B"; radius: 8
                border.color: subUrlField.activeFocus ? "#00FF9D" : "#22222E"; border.width: 1
            }
        }

        // Server List Section
        Text {
            text: "AVAILABLE NODES"
            color: "#5A5A6A"
            font.pixelSize: 11
            font.bold: true
            font.letterSpacing: 1
            Layout.topMargin: 5
        }

        ListView {
            Layout.fillWidth: true
            Layout.fillHeight: true
            clip: true
            model: serverList

            delegate: Rectangle {
                width: ListView.view.width
                height: 60
                color: "#14141B"
                radius: 10
                border.color: "#22222E"
                border.width: 1
                Layout.bottomMargin: 8

                RowLayout {
                    anchors.fill: parent
                    anchors.margins: 14
                    spacing: 12

                    Rectangle {
                        width: 8; height: 8
                        radius: 4
                        color: modelData.latency < 100 ? "#00FF9D" : (modelData.latency < 300 ? "#FFCC00" : "#FF3366")
                    }

                    ColumnLayout {
                        Layout.fillWidth: true
                        spacing: 2
                        Text { text: modelData.name; color: "#FFFFFF"; font.bold: true; font.pixelSize: 13; elide: Text.ElideRight; Layout.fillWidth: true }
                        Text { text: modelData.protocol.toUpperCase() + " • " + modelData.address; color: "#5A5A6A"; font.pixelSize: 10 }
                    }
                    Text {
                        text: modelData.latency == 9999 ? "N/A" : (modelData.latency == 0 ? "-" : modelData.latency + "ms")
                        color: "#5A5A6A"; font.bold: true; font.pixelSize: 11
                    }
                }
                MouseArea {
                    anchors.fill: parent
                    cursorShape: Qt.PointingHandCursor
                    onClicked: { prevDownload = 0; prevUpload = 0; prevTime = Date.now(); vpnController.connectToServer(modelData.uri) }
                }
            }
        }
    }
}
