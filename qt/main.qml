import QtQuick
import QtQuick.Controls
import QtQuick.Controls.Material
import QtQuick.Layouts

Window {
    id: mainWindow
    width: 400
    height: 750
    visible: true
    title: "Kingo VPN"
    color: "#101016" // Hiddify Dark Background

    // Fix Qt6 Syntax: Explicitly declare parameters
    onClosing: (close) => {
        close.accepted = false
        hide()
        trayIcon.showMessage("Kingo VPN", "Application minimized to tray.")
    }

    property real prevDownload: 0
    property real prevUpload: 0
    property real prevTime: 0

    Timer {
        interval: 1000
        running: vpnController.status == "connected"
        repeat: true
        onTriggered: vpnController.getTraffic()
    }

    function formatSpeed(bytes) {
        if (!bytes || bytes < 0) return "0 B/s"
        if (bytes < 1024) return bytes.toFixed(0) + " B/s"
        if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + " KB/s"
        return (bytes / (1024 * 1024)).toFixed(2) + " MB/s"
    }

    Connections {
        target: trayIcon
        function onActivateRequested() {
            mainWindow.show()
            mainWindow.raise()
            mainWindow.requestActivate()
        }
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
        spacing: 20

        // Header
        Label {
            text: "Kingo VPN"
            color: "#FFFFFF"
            font.pixelSize: 28
            font.bold: true
            Layout.alignment: Qt.AlignHCenter
            Layout.topMargin: 20
        }

        // Status Card (Hiddify Style)
        Rectangle {
            Layout.fillWidth: true
            height: 180
            color: "#1E1E2A"
            radius: 24

            ColumnLayout {
                anchors.fill: parent
                anchors.margins: 16
                spacing: 8

                Label {
                    text: vpnController.status.toUpperCase()
                    color: vpnController.status == "connected" ? "#3DDC84" : (vpnController.status == "connecting" ? "#FFCC00" : "#FF5252")
                    font.pixelSize: 20
                    font.bold: true
                    Layout.alignment: Qt.AlignHCenter
                }

                RowLayout {
                    Layout.alignment: Qt.AlignHCenter
                    spacing: 40
                    visible: vpnController.status == "connected"

                    ColumnLayout {
                        spacing: 2
                        Label {
                            text: "Download"
                            color: "#B0B0B0"
                            font.pixelSize: 12
                            Layout.alignment: Qt.AlignHCenter
                        }
                        Label {
                            id: downloadSpeedText
                            text: "0 B/s"
                            color: "#3DDC84"
                            font.pixelSize: 18
                            font.bold: true
                            Layout.alignment: Qt.AlignHCenter
                        }
                    }

                    ColumnLayout {
                        spacing: 2
                        Label {
                            text: "Upload"
                            color: "#B0B0B0"
                            font.pixelSize: 12
                            Layout.alignment: Qt.AlignHCenter
                        }
                        Label {
                            id: uploadSpeedText
                            text: "0 B/s"
                            color: "#FFCC00"
                            font.pixelSize: 18
                            font.bold: true
                            Layout.alignment: Qt.AlignHCenter
                        }
                    }
                }
            }
        }

        // Big Connect Button (Hiddify Style)
        Button {
            Layout.fillWidth: true
            height: 60
            enabled: vpnController.status != "connecting"
            
            background: Rectangle {
                color: vpnController.status == "connected" ? "#FF5252" : "#3DDC84"
                radius: 30
            }
            
            contentItem: Text {
                text: vpnController.status == "connecting" ? "Connecting..." : (vpnController.status == "connected" ? "DISCONNECT" : "CONNECT")
                color: "#101016"
                font.pixelSize: 18
                font.bold: true
                horizontalAlignment: Text.AlignHCenter
                verticalAlignment: Text.AlignVCenter
            }

            onClicked: {
                if (vpnController.status == "connected") {
                    vpnController.disconnectVpn()
                } else {
                    prevDownload = 0
                    prevUpload = 0
                    prevTime = Date.now()
                    vpnController.autoConnect()
                }
            }
        }

        // Error Box
        Rectangle {
            id: errorBox
            Layout.fillWidth: true
            height: 40
            color: "#FF5252"
            radius: 12
            visible: false
            
            Label {
                id: errorText
                anchors.centerIn: parent
                color: "#FFFFFF"
                font.bold: true
            }
            
            Timer {
                id: errorTimer
                interval: 4000
                onTriggered: errorBox.visible = false
            }
        }

        // Subscription Input
        TextField {
            id: subUrlField
            Layout.fillWidth: true
            height: 50
            placeholderText: "Enter Subscription URL..."
            color: "#FFFFFF"
            font.pixelSize: 14
            text: "https://raw.githubusercontent.com/MhdiTaheri/VpnHub/main/sub"
            
            background: Rectangle {
                color: "#1E1E2A"
                radius: 16
                border.color: subUrlField.activeFocus ? "#3DDC84" : "transparent"
                border.width: 2
            }
        }

        Button {
            Layout.fillWidth: true
            height: 44
            text: "Update Subscription"
            
            background: Rectangle {
                color: "#2A2A3A"
                radius: 16
            }
            contentItem: Text {
                text: parent.text
                color: "#FFFFFF"
                font.bold: true
                font.pixelSize: 14
                horizontalAlignment: Text.AlignHCenter
                verticalAlignment: Text.AlignVCenter
            }
            
            onClicked: {
                if (subUrlField.text.length > 0) {
                    vpnController.addSubscription(subUrlField.text)
                }
            }
        }

        // Servers List
        Label {
            text: "Available Servers"
            color: "#B0B0B0"
            font.pixelSize: 14
            font.bold: true
            Layout.topMargin: 10
        }

        ListView {
            Layout.fillWidth: true
            Layout.fillHeight: true
            clip: true
            model: vpnController.servers

            delegate: Rectangle {
                width: ListView.view.width
                height: 65
                color: "#1E1E2A"
                radius: 16
                Layout.bottomMargin: 10

                RowLayout {
                    anchors.fill: parent
                    anchors.margins: 16
                    spacing: 12

                    ColumnLayout {
                        Layout.fillWidth: true
                        spacing: 2
                        
                        Label {
                            text: modelData.name
                            color: "#FFFFFF"
                            font.bold: true
                            font.pixelSize: 14
                            elide: Text.ElideRight
                            Layout.fillWidth: true
                        }
                        Label {
                            text: modelData.protocol + " - " + modelData.address
                            color: "#B0B0B0"
                            font.pixelSize: 11
                        }
                    }

                    Label {
                        text: modelData.latency == 9999 ? "N/A" : (modelData.latency == 0 ? "-" : modelData.latency + "ms")
                        color: modelData.latency < 100 ? "#3DDC84" : (modelData.latency < 300 ? "#FFCC00" : "#FF5252")
                        font.bold: true
                        font.pixelSize: 12
                    }
                }

                MouseArea {
                    anchors.fill: parent
                    cursorShape: Qt.PointingHandCursor
                    onClicked: {
                        prevDownload = 0
                        prevUpload = 0
                        prevTime = Date.now()
                        vpnController.connectToServer(modelData.uri)
                    }
                }
            }
        }
    }
}
