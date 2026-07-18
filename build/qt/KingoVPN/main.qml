import QtQuick
import QtQuick.Controls

Window {
    width: 400
    height: 600
    visible: true
    title: "Kingo Linux VPN"
    color: "#1e1e2e"

    Column {
        anchors.centerIn: parent
        spacing: 20

        Text {
            anchors.horizontalCenter: parent.horizontalCenter
            text: "Status: " + (vpnController.connected ? "Connected" : "Disconnected")
            color: vpnController.connected ? "#a6e3a1" : "#f38ba8"
            font.pointSize: 24
            font.bold: true
        }

        Button {
            anchors.horizontalCenter: parent.horizontalCenter
            text: vpnController.connected ? "Disconnect" : "Connect (Mock Config)"
            onClicked: {
                if (vpnController.connected) {
                    vpnController.disconnectVpn()
                } else {
                    // برای تست اولیه یک مسیر فیک ارسال می‌کنیم
                    vpnController.connectVpn("/tmp/dummy_config.json")
                }
            }
        }
    }
}
