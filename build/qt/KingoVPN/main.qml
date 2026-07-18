import QtQuick
import QtQuick.Controls

Window {
    width: 400
    height: 600
    visible: true
    title: "Kingo Linux VPN"
    Rectangle {
        anchors.fill: parent
        color: "#1e1e2e"
        Text {
            anchors.centerIn: parent
            text: "Kingo VPN\nBackend Running..."
            color: "#cdd6f4"
            font.pointSize: 24
            horizontalAlignment: Text.AlignHCenter
        }
    }
}
