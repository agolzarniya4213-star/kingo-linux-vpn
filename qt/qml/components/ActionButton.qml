import QtQuick
import QtQuick.Controls

Button {
    property bool secondary: false
    padding: 10

    background: Rectangle {
        radius: 14
        color: secondary ? "#223244" : "#2f7cf6"
        border.color: "#31445c"
        border.width: 1
    }

    contentItem: Text {
        text: parent.text
        color: "white"
        font.bold: true
        horizontalAlignment: Text.AlignHCenter
        verticalAlignment: Text.AlignVCenter
    }
}
