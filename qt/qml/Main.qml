import QtQuick 2.15
import QtQuick.Controls 2.15
import QtQuick.Controls.Material 2.15
import QtQuick.Layouts 1.15

ApplicationWindow {
    id: window
    visible: true
    width: 420
    height: 740
    color: "#0f172a"
    title: "Kingo VPN"
    Material.theme: Material.Dark
    Material.accent: Material.Blue

    // FIX 1: Removed redundant close button. Rely only on Window Manager.
    Rectangle {
        id: topBar
        width: parent.width
        height: 40
        color: "#0f172a"
        
        Text {
            text: "KINGO VPN"
            color: "#ffffff"
            font.pixelSize: 16
            font.bold: true
            font.letterSpacing: 1
            anchors.left: parent.left
            anchors.leftMargin: 15
            anchors.verticalCenter: parent.verticalCenter
        }
    }

    Item {
        anchors.fill: parent
        anchors.topMargin: 40

        ColumnLayout {
            anchors.fill: parent
            anchors.margins: 25
            spacing: 20

            // Free Plan Card with Live Stats
            Rectangle {
                Layout.fillWidth: true
                Layout.preferredHeight: 85
                radius: 16
                color: "#1e293b"
                
                RowLayout {
                    anchors.fill: parent
                    anchors.margins: 15
                    spacing: 15
                    
                    Rectangle {
                        width: 40
                        height: 40
                        radius: 20
                        color: "#3b82f6"
                        Text { 
                            text: "K"
                            color: "white"
                            font.bold: true
                            font.pixelSize: 20
                            anchors.centerIn: parent
                        }
                    }
                    Column {
                        Layout.fillWidth: true
                        spacing: 4
                        Text { 
                            text: "Free Plan" 
                            color: "#94a8b8"
                            font.pixelSize: 12
                        }
                        Text { 
                            text: "⬇ " + appCore.rxData + "   ⬆ " + appCore.txData
                            color: "#e2e8f0"
                            font.bold: true
                            font.pixelSize: 13
                        }
                        Text { 
                            text: "PING: " + appCore.pingData
                            color: "#4ade80"
                            font.pixelSize: 12
                        }
                    }
                }
            }

            Item {
                Layout.fillHeight: true
            }

            // Premium Connection Area
            Item {
                Layout.alignment: Qt.AlignHCenter
                width: 220
                height: 220

                Rectangle {
                    id: outerRing
                    anchors.fill: parent
                    radius: width / 2
                    color: "transparent"
                    border.color: appCore.isConnected ? "#4ade80" : (connectMa.pressed ? "#3b82f6" : "#334155")
                    border.width: 3
                    scale: appCore.isConnected ? 1.0 : 1.05
                    opacity: appCore.isConnected ? 0.8 : 0.5
                    Behavior on scale { NumberAnimation { duration: 300 } }
                    Behavior on opacity { NumberAnimation { duration: 300 } }
                }

                Rectangle {
                    id: innerCircle
                    anchors.centerIn: parent
                    width: 180
                    height: 180
                    radius: 90
                    color: appCore.isConnected ? "#166534" : "#1e293b"
                    border.color: appCore.isConnected ? "#4ade80" : "#334155"
                    border.width: 2
                    scale: connectMa.pressed ? 0.95 : 1.0
                    
                    Behavior on scale { NumberAnimation { duration: 150 } }
                    Behavior on color { ColorAnimation { duration: 400 } }

                    Text {
                        text: appCore.isConnected ? "CONNECTED" : "CONNECT"
                        color: "white"
                        font.bold: true
                        font.pixelSize: 22
                        font.letterSpacing: 2
                        anchors.centerIn: parent
                    }
                }

                MouseArea {
                    id: connectMa
                    anchors.fill: parent
                    cursorShape: Qt.PointingHandCursor
                    onClicked: appCore.toggleConnection()
                }
            }

            Text {
                text: appCore.statusText
                color: appCore.isConnected ? "#4ade80" : "#94a3b8"
                font.pixelSize: 16
                font.bold: true
                Layout.alignment: Qt.AlignHCenter
                horizontalAlignment: Text.AlignHCenter
            }

            Item {
                Layout.fillHeight: true
            }

            // Server Selection Card
            Rectangle {
                Layout.fillWidth: true
                Layout.preferredHeight: 75
                radius: 16
                color: "#1e293b"
                
                RowLayout {
                    anchors.fill: parent
                    anchors.margins: 15
                    spacing: 15

                    Rectangle { 
                        width: 40
                        height: 28
                        radius: 4
                        color: "#3b82f6"
                    }
                    
                    Column {
                        Layout.fillWidth: true
                        spacing: 2
                        Text { 
                            text: appCore.selectedServer
                            color: "white"
                            font.bold: true
                            font.pixelSize: 15
                        }
                        Text { 
                            text: "UDP Protocol"
                            color: "#64748b"
                            font.pixelSize: 12
                        }
                    }

                    Text { 
                        text: "›"
                        font.pixelSize: 30
                        color: "#64748b"
                    }
                }

                MouseArea {
                    anchors.fill: parent
                    cursorShape: Qt.PointingHandCursor
                    onClicked: serverPopup.open()
                }
            }

            // Bottom Navigation
            RowLayout {
                Layout.fillWidth: true
                Layout.preferredHeight: 55
                spacing: 10
                
                Rectangle {
                    Layout.fillWidth: true
                    Layout.fillHeight: true
                    radius: 14
                    color: "#3b82f6"
                    Text { 
                        text: "Home"
                        color: "white"
                        font.pixelSize: 14
                        font.bold: true
                        anchors.centerIn: parent
                    }
                }
                Rectangle {
                    Layout.fillWidth: true
                    Layout.fillHeight: true
                    radius: 14
                    color: "#1e293b"
                    Text { 
                        text: "Servers"
                        color: "#94a3b8"
                        font.pixelSize: 14
                        anchors.centerIn: parent
                    }
                    MouseArea { 
                        anchors.fill: parent
                        cursorShape: Qt.PointingHandCursor
                        onClicked: serverPopup.open()
                    }
                }
                Rectangle {
                    Layout.fillWidth: true
                    Layout.fillHeight: true
                    radius: 14
                    color: "#1e293b"
                    Text { 
                        text: "Settings"
                        color: "#94a3b8"
                        font.pixelSize: 14
                        anchors.centerIn: parent
                    }
                    MouseArea { 
                        anchors.fill: parent
                        cursorShape: Qt.PointingHandCursor
                        onClicked: settingsPopup.open()
                    }
                }
            }
        }
    }

    // FIX 2 & 3: Hiddify-Style Server Bottom Sheet (Strict Layouts + Fix setServer)
    Popup {
        id: serverPopup
        y: parent.height - height
        width: parent.width
        height: parent.height * 0.6
        padding: 20
        modal: true
        closePolicy: Popup.CloseOnEscape | Popup.CloseOnPressOutside
        background: Rectangle { color: "#1e293b"; radius: 16 }

        Column {
            anchors.fill: parent
            spacing: 15
            
            Text {
                text: "Select Server"
                color: "white"
                font.bold: true
                font.pixelSize: 18
            }
            
            Rectangle { 
                width: parent.width
                height: 1
                color: "#334155" 
            }
            
            Repeater {
                model: ["US - New York #1", "US - Los Angeles #2", "DE - Frankfurt #1", "NL - Amsterdam #1", "UK - London #1", "JP - Tokyo #1"]
                delegate: Rectangle {
                    width: serverPopup.width - 40
                    height: 50
                    radius: 10
                    color: appCore.selectedServer === modelData ? "#334155" : "transparent"
                    
                    RowLayout {
                        anchors.fill: parent
                        anchors.margins: 10
                        spacing: 15
                        
                        Rectangle { 
                            width: 30
                            height: 20
                            radius: 3
                            color: "#3b82f6" 
                        }
                        
                        Text { 
                            text: modelData
                            color: "white"
                            font.pixelSize: 15
                            Layout.fillWidth: true
                        }
                        
                        // FIX: Removed invalid anchors. Layout.alignment handles it automatically.
                        Text {
                            text: appCore.selectedServer === modelData ? "✓" : ""
                            color: "#4ade80"
                            font.pixelSize: 18
                        }
                    }
                    
                    MouseArea {
                        anchors.fill: parent
                        cursorShape: Qt.PointingHandCursor
                        onClicked: {
                            // FIX: Use QML property assignment instead of C++ setter function
                            appCore.selectedServer = modelData;
                            serverPopup.close();
                        }
                    }
                }
            }
        }
    }

    // Settings Placeholder
    Popup {
        id: settingsPopup
        anchors.centerIn: parent
        width: parent.width - 50
        height: 200
        modal: true
        closePolicy: Popup.CloseOnEscape | Popup.CloseOnPressOutside
        background: Rectangle { color: "#1e293b"; radius: 16 }
        
        Column {
            anchors.centerIn: parent
            spacing: 15
            Text {
                text: "Settings"
                color: "white"
                font.bold: true
                font.pixelSize: 18
                anchors.horizontalCenter: parent.horizontalCenter
            }
            Rectangle { 
                width: 200
                height: 1
                color: "#334155" 
                anchors.horizontalCenter: parent.horizontalCenter
            }
            Text {
                text: "Advanced settings coming soon..."
                color: "#64748b"
                font.pixelSize: 14
                anchors.horizontalCenter: parent.horizontalCenter
            }
        }
    }
}
