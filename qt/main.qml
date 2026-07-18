import QtQuick
import QtQuick.Controls

Window {
    width: 400
    height: 600
    visible: true
    title: "Kingo Linux VPN"
    color: "#1e1e2e"

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

        Button {
            text: vpnController.connected ? "Disconnect" : "Connect (Mock Config)"
            onClicked: {
                if (vpnController.connected) {
                    vpnController.disconnectVpn()
                } else {
                    vpnController.connectVpn("/tmp/dummy_config.json")
                }
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
            height: 400
            clip: true
            model: vpnController.servers

            delegate: Rectangle {
                width: parent.width
                height: 50
                color: "#313244"
                radius: 5

                Column {
                    anchors.centerIn: parent
                    Text {
                        text: modelData.name
                        color: "#cdd6f4"
                        font.bold: true
                    }
                    Text {
                        text: modelData.protocol + " - " + modelData.address
                        color: "#a6adc8"
                        font.pointSize: 8
                    }
                }
            }
        }
    }
}
