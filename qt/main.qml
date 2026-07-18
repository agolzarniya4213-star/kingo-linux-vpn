import QtQuick
import QtQuick.Controls

Window {
    width: 400
    height: 700
    visible: true
    title: "Kingo Linux VPN"
    color: "#1e1e2e"

    // تایمر برای دریافت ترافیک هر 1 ثانیه
    Timer {
        interval: 1000
        running: vpnController.connected
        repeat: true
        onTriggered: vpnController.getTraffic()
    }

    // تابع کمکی برای فرمت کردن سرعت
    function formatSpeed(bytes) {
        if (bytes < 1024) return bytes + " B/s"
        if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + " KB/s"
        return (bytes / (1024 * 1024)).toFixed(2) + " MB/s"
    }

    Column {
        anchors.fill: parent
        anchors.margins: 20
        spacing: 15

        Text {
            text: "Status: " + (vpnController.connected ? "Connected" : "Disconnected")
            color: vpnController.connected ? "#a6e3a1" : "#f38ba8"
            font.pointSize: 16
            font.bold: true
        }

        // نمایش ترافیک
        Row {
            width: parent.width
            spacing: 20
            visible: vpnController.connected

            Column {
                Text {
                    text: "Download"
                    color: "#a6adc8"
                    font.pointSize: 10
                }
                Text {
                    text: formatSpeed(vpnController.downloadSpeed)
                    color: "#89b4fa"
                    font.pointSize: 18
                    font.bold: true
                }
            }

            Column {
                Text {
                    text: "Upload"
                    color: "#a6adc8"
                    font.pointSize: 10
                }
                Text {
                    text: formatSpeed(vpnController.uploadSpeed)
                    color: "#f9e2af"
                    font.pointSize: 18
                    font.bold: true
                }
            }
        }

        Button {
            text: vpnController.connected ? "Disconnect" : "Select a server to connect"
            enabled: vpnController.connected
            onClicked: vpnController.disconnectVpn()
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
                text: "Update Subscription"
                onClicked: {
                    if (subUrlField.text.length > 0) {
                        vpnController.addSubscription(subUrlField.text)
                    }
                }
            }

            Button {
                text: "Test Latency"
                onClicked: vpnController.testLatency()
            }
        }

        Text {
            text: "Available Servers (Click to connect):"
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
                    onClicked: {
                        vpnController.connectToServer(modelData.uri)
                    }
                }
            }
        }
    }
}
