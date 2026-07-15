import QtQuick
import QtQuick.Controls
import QtQuick.Layouts

Frame {
    property string title: ""
    property string subtitle: ""
    padding: 16

    background: Rectangle {
        radius: 18
        color: "#121a25"
        border.color: "#223042"
    }

    ColumnLayout {
        anchors.fill: parent
        spacing: 4

        Label { text: title; font.pixelSize: 18; font.bold: true }
        Label { text: subtitle; opacity: 0.72; wrapMode: Text.WordWrap }
    }
}
