import QtQuick
import QtQuick.Controls
import QtQuick.Layouts

Frame {
    id: rowRoot
    property string serverName: ""
    property string serverProtocol: ""
    property int serverPing: -1
    property string serverGroup: "servers"
    signal connectRequested()

    padding: 12

    background: Rectangle {
        radius: 14
        color: "#0f1620"
        border.color: "#223042"
    }

    RowLayout {
        anchors.fill: parent
        spacing: 12

        ColumnLayout {
            Layout.fillWidth: true
            spacing: 4

            Label { text: serverName; font.bold: true; font.pixelSize: 15 }
            Label { text: serverProtocol + " · ping " + serverPing + " ms"; opacity: 0.72 }
            Label { text: serverGroup; opacity: 0.55 }
        }

        Button {
            text: "Connect"
            onClicked: rowRoot.connectRequested()
        }
    }
}
