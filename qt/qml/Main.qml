import QtQuick
import QtQuick.Controls
import QtQuick.Layouts
import KingoVpn 1.0

ApplicationWindow {
    width: 1280
    height: 820
    visible: true
    title: "Kingo Linux VPN"

    header: ToolBar {
        RowLayout {
            anchors.fill: parent
            anchors.margins: 12
            spacing: 10

            Label {
                text: "Kingo Linux VPN"
                font.pixelSize: 20
                font.bold: true
            }

            Item { Layout.fillWidth: true }

            Label {
                text: "State: " + AppController.state
                opacity: 0.8
            }

            Button { text: "Refresh"; onClicked: { AppController.refresh(); AppController.reloadStatus() } }
            Button { text: "Reload"; onClicked: AppController.reloadStatus() }
            Button { text: "Stop"; onClicked: AppController.stopEngine() }
        }
    }

    SplitView {
        anchors.fill: parent
        orientation: Qt.Horizontal

        Frame {
            SplitView.preferredWidth: 760
            ColumnLayout {
                anchors.fill: parent
                spacing: 12

                InfoCard {
                    title: "Servers"
                    subtitle: "Servers received from the daemon cache"
                }

                ListView {
                    Layout.fillWidth: true
                    Layout.fillHeight: true
                    clip: true
                    model: ServerListModel
                    delegate: ServerRow {
                        serverName: model.name
                        serverProtocol: model.protocol
                        serverPing: model.pingMs
                        serverGroup: model.group
                        onConnectRequested: AppController.startServer(ServerListModel.get(index))
                    }
                }
            }
        }

        Frame {
            SplitView.preferredWidth: 420
            ColumnLayout {
                anchors.fill: parent
                spacing: 12

                InfoCard {
                    title: "Connection"
                    subtitle: "Daemon socket: " + AppController.socketPath
                }

                TextArea {
                    id: jsonInput
                    Layout.fillWidth: true
                    Layout.preferredHeight: 170
                    placeholderText: "Paste server JSON here"
                }

                RowLayout {
                    Layout.fillWidth: true
                    spacing: 10

                    ActionButton {
                        text: "Connect"
                        onClicked: AppController.startFromJson(jsonInput.text)
                    }

                    ActionButton {
                        text: "Reload status"
                        secondary: true
                        onClicked: AppController.reloadStatus()
                    }
                }

                InfoCard {
                    title: "Status"
                    subtitle: AppController.lastError.length > 0 ? AppController.lastError : "The UI mirrors daemon responses directly."
                }
            }
        }
    }

    Component.onCompleted: AppController.reloadStatus()
}
