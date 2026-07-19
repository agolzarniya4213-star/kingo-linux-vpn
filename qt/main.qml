import QtQuick
import QtQuick.Controls
import QtQuick.Controls.Material
import QtQuick.Layouts

Window {
    id: mainWindow
    width: 400
    height: 700
    visible: true
    title: "Kingo VPN"
    color: "#1e1e2e" // Catppuccin Dark background like Hiddify

    // Fix Qt6 Syntax: Explicitly declare parameters
    onClosing: (close) => {
        close.accepted = false
        hide()
        trayIcon.showMessage("Kingo VPN", "Application minimized to tray.")
    }

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

    // Real-time speed calculation
    property real prevDownload: 0
    property real prevUpload: 0
    property real prevTime: 0

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
        // Fix Qt6 Syntax: Explicitly declare parameters
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
            color: "#cdd6f4"
            font.pixelSize: 28
            font.bold: true
            Layout.alignment: Qt.AlignHCenter
        }

        // Status Card
        Rectangle {
            Layout.fillWidth: true
            height: 160
            color: "#313244"
            radius: 16

            ColumnLayout {
                anchors.fill: parent
                anchors.margins: 16
                spacing: 8

                Label {
                    text: vpnController.status.toUpperCase()
                    color: vpnController.status == "connected" ? "#a6e3a1" : (vpnController.status == "connecting" ? "#f9e2af" : "#f38ba8")
                    font.pixelSize: 20
                    font.bold: true
                    Layout.alignment: Qt.AlignHCenter
                }

                RowLayout {
                    Layout.alignment: Qt.AlignHCenter
                    spacing: 24
                    visible: vpnController.status == "connected"

                    ColumnLayout {
                        spacing: 2
                        Label {
                            text: "Download"
                            color: "#a6adc8"
                            font.pixelSize: 12
                            Layout.alignment: Qt.AlignHCenter
                        }
                        Label {
                            id: downloadSpeedText
                            text: "0 B/s"
                            color: "#89b4fa"
                            font.pixelSize: 18
                            font.bold: true
                            Layout.alignment: Qt.AlignHCenter
                        }
                    }

                    ColumnLayout {
                        spacing: 2
                        Label {
                            text: "Upload"
                            color: "#a6adc8"
                            font.pixelSize: 12
                            Layout.alignment: Qt.AlignHCenter
                        }
                        Label {
                            id: uploadSpeedText
                            text: "0 B/s"
                            color: "#f9e2af"
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
            height: 56
            text: vpnController.status == "connecting" ? "Connecting..." : (vpnController.status == "connected" ? "Disconnect" : "Auto Connect")
            enabled: vpnController.status != "connecting"
            
            Material.background: vpnController.status == "connected" ? "#f38ba8" : "#a6e3a1"
            Material.foreground: "#1e1e2e"
            font.pixelSize: 18
            font.bold: true
            radius: 16

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
            color: "#f38ba8"
            radius: 8
            visible: false
            
            Label {
                id: errorText
                anchors.centerIn: parent
                color: "#1e1e2e"
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
            height: 48
            placeholderText: "Enter Subscription URL..."
            color: "#cdd6f4"
            font.pixelSize: 14
            
            background: Rectangle {
                color: "#313244"
                radius: 12
                border.color: subUrlField.activeFocus ? "#89b4fa" : "transparent"
                border.width: 2
            }
        }

        Button {
            Layout.fillWidth: true
            height: 44
            text: "Update Subscription"
            Material.background: "#313244"
            Material.foreground: "#cdd6f4"
            font.pixelSize: 14
            radius: 12
            
            onClicked: {
                if (subUrlField.text.length > 0) {
                    vpnController.addSubscription(subUrlField.text)
                }
            }
        }

        // Servers List
        Label {
            text: "Available Servers"
            color: "#a6adc8"
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
                height: 60
                color: "transparent"
                
                Rectangle {
                    anchors.fill: parent
                    anchors.margins: 2
                    color: "#313244"
                    radius: 12

                    RowLayout {
                        anchors.fill: parent
                        anchors.margins: 12
                        spacing: 12

                        ColumnLayout {
                            Layout.fillWidth: true
                            spacing: 2
                            
                            Label {
                                text: modelData.name
                                color: "#cdd6f4"
                                font.bold: true
                                font.pixelSize: 14
                                elide: Text.ElideRight
                                Layout.fillWidth: true
                            }
                            Label {
                                text: modelData.protocol + " - " + modelData.address
                                color: "#a6adc8"
                                font.pixelSize: 11
                            }
                        }

                        Label {
                            text: modelData.latency == 9999 ? "N/A" : (modelData.latency == 0 ? "-" : modelData.latency + "ms")
                            color: modelData.latency < 100 ? "#a6e3a1" : (modelData.latency < 300 ? "#f9e2af" : "#f38ba8")
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
}
