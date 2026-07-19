import QtQuick
import QtQuick.Controls

Window {
    id: mainWindow
    width: 400
    height: 700
    visible: true
    title: "Kingo Linux VPN"
    color: "#1e1e2e"

    // FIX BUG-024: Variables to calculate real speed (Delta over time)
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
        if (bytes < 1024) return bytes.toFixed(0) + " B/s"
        if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + " KB/s"
        return (bytes / (1024 * 1024)).toFixed(2) + " MB/s"
    }

    onClosing: {
        close.accepted = false
        hide()
        trayIcon.showMessage("Kingo VPN", "Application minimized to tray.")
    }

    Connections {
        target: trayIcon
        function onActivateRequested() {
            mainWindow.show()
            mainWindow.raise()
            mainWindow.requestActivate()
        }
        function onConnectRequested() {
            mainWindow.show()
            mainWindow.raise()
        }
        function onDisconnectRequested() {
            vpnController.disconnectVpn()
        }
        function onQuitRequested() {
            Qt.quit()
        }
    }

    Connections {
        target: vpnController
        function onErrorOccurred(error) {
            errorText.text = error
            errorText.visible = true
        }
        // FIX BUG-024: Calculate speed when traffic changes
        function onTrafficChanged() {
            var currTime = Date.now()
            var dt = (currTime - prevTime) / 1000.0
            
            if (dt > 0) {
                var dlSpeed = (vpnController.downloadSpeed - prevDownload) / dt
                var ulSpeed = (vpnController.uploadSpeed - prevUpload) / dt
                downloadSpeedText.text = formatSpeed(dlSpeed)
                uploadSpeedText.text = formatSpeed(ulSpeed)
            }
            
            prevDownload = vpnController.downloadSpeed
            prevUpload = vpnController.uploadSpeed
            prevTime = currTime
        }
    }

    Column {
        anchors.fill: parent
        anchors.margins: 20
        spacing: 15

        Text {
            text: "Status: " + vpnController.status
            color: vpnController.status == "connected" ? "#a6e3a1" : (vpnController.status == "connecting" ? "#f9e2af" : "#f38ba8")
            font.pointSize: 16
            font.bold: true
        }

        Row {
            width: parent.width
            spacing: 20
            visible: vpnController.status == "connected"

            Column {
                Text { text: "Download"; color: "#a6adc8"; font.pointSize: 10 }
                Text {
                    id: downloadSpeedText
                    text: "0 B/s"
                    color: "#89b4fa"; font.pointSize: 18; font.bold: true
                }
            }

            Column {
                Text { text: "Upload"; color: "#a6adc8"; font.pointSize: 10 }
                Text {
                    id: uploadSpeedText
                    text: "0 B/s"
                    color: "#f9e2af"; font.pointSize: 18; font.bold: true
                }
            }
        }

        Row {
            width: parent.width
            spacing: 10

            Button {
                text: vpnController.status == "connecting" ? "Connecting..." : "Auto Connect"
                enabled: vpnController.status != "connecting"
                highlighted: true
                onClicked: {
                    // Reset deltas on connect
                    prevDownload = 0
                    prevUpload = 0
                    prevTime = Date.now()
                    vpnController.autoConnect()
                }
            }

            Button {
                text: "Disconnect"
                enabled: vpnController.status == "connected"
                onClicked: vpnController.disconnectVpn()
            }
        }

        Rectangle {
            id: errorBox
            width: parent.width
            height: errorText.visible ? 40 : 0
            color: "#f38ba8"
            radius: 5
            visible: errorText.visible
            Text {
                id: errorText
                anchors.centerIn: parent
                color: "#11111b"
                font.bold: true
                visible: false
                onTextChanged: {
                    if (text !== "") {
                        visible = true
                        errorTimer.start()
                    }
                }
                Timer {
                    id: errorTimer
                    interval: 5000
                    onTriggered: {
                        errorText.visible = false
                        errorText.text = ""
                    }
                }
            }
        }

        TextField {
            id: subUrlField
            width: parent.width
            height: 40
            placeholderText: "Enter Subscription URL..."
            color: "#cdd6f4"
            background: Rectangle { color: "#313244"; radius: 5 }
        }

        Row {
            width: parent.width
            spacing: 10
            Button {
                text: "Update Sub"
                onClicked: {
                    if (subUrlField.text.length > 0) {
                        vpnController.addSubscription(subUrlField.text)
                    }
                }
            }
            Button {
                text: "Test Ping"
                onClicked: vpnController.testLatency()
            }
        }

        Text {
            text: "Available Servers:"
            color: "#cdd6f4"
            font.pointSize: 14
            topPadding: 10
        }

        ListView {
            width: parent.width
            height: 350
            clip: true
            model: vpnController.servers

            delegate: Rectangle {
                width: parent.width
                height: 50
                color: "#313244"
                radius: 5

                Row {
                    anchors.fill: parent
                    anchors.margins: 10
                    spacing: 10

                    Column {
                        width: parent.width - 60
                        anchors.verticalCenter: parent.verticalCenter
                        Text {
                            text: modelData.name
                            color: "#cdd6f4"
                            font.bold: true
                            elide: Text.ElideRight
                            width: parent.width
                        }
                        Text {
                            text: modelData.protocol + " - " + modelData.address
                            color: "#a6adc8"
                            font.pointSize: 8
                        }
                    }

                    Text {
                        width: 50
                        anchors.verticalCenter: parent.verticalCenter
                        horizontalAlignment: Text.AlignRight
                        text: modelData.latency == 9999 ? "N/A" : (modelData.latency == 0 ? "-" : modelData.latency + "ms")
                        color: modelData.latency < 100 ? "#a6e3a1" : (modelData.latency < 300 ? "#f9e2af" : "#f38ba8")
                        font.bold: true
                    }
                }

                MouseArea {
                    anchors.fill: parent
                    cursorShape: Qt.PointingHandCursor
                    onClicked: vpnController.connectToServer(modelData.uri)
                }
            }
        }
    }
}
