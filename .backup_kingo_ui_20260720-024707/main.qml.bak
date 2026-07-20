import QtQuick
import QtQuick.Controls
import QtQuick.Layouts

ApplicationWindow {
    id: mainWindow
    width: 520
    height: 920
    minimumWidth: 460
    minimumHeight: 820
    visible: true
    title: "Kingo VPN"
    color: bgColor

    property string connStatus: vpnController ? vpnController.status : "disconnected"
    property int currentTab: 1
    property string selectedServerID: ""
    property bool showSettings: false
    property bool autoScanActive: false
    property int autoScanProgress: 0

    property string subscriptionUrl: "https://raw.githubusercontent.com/MhdiTaheri/VpnHub/main/sub"

    property color bgColor: "#070D18"
    property color surfaceColor: "#1A2636"
    property color surface2Color: "#0E1A2C"
    property color borderColor: "#243245"
    property color accentColor: "#18CFFF"
    property color accent2Color: "#0FA8D9"
    property color textColor: "#FFFFFF"
    property color subTextColor: "#8899AA"
    property color successColor: "#4CAF50"
    property color warningColor: "#FFC107"
    property color errorColor: "#F44336"

    function safeValue(v, fallback) {
        if (v === undefined || v === null || v === "") return fallback
        return v
    }

    function serverKey(server, idx) {
        if (!server) return String(idx)
        return safeValue(server.id, safeValue(server.uri, String(idx)))
    }

    function serverName(server) {
        return safeValue(server.name, safeValue(server.title, "Unknown Server"))
    }

    function serverProtocol(server) {
        return String(safeValue(server.protocol, safeValue(server.type, "VLESS"))).toUpperCase()
    }

    function serverAddress(server) {
        return safeValue(server.address, safeValue(server.host, safeValue(server.domain, "")))
    }

    function serverLatency(server) {
        var v = safeValue(server.latency, 9999)
        return Number(v)
    }

    function getPingColor(latency) {
        if (latency === 9999 || latency === 0) return subTextColor
        if (latency < 150) return successColor
        if (latency < 300) return warningColor
        return errorColor
    }

    function formatSpeed(bytesPerSecond) {
        var v = Number(bytesPerSecond || 0)
        if (v < 0) v = 0
        if (v < 1024) return v.toFixed(1) + " B/s"
        if (v < 1048576) return (v / 1024).toFixed(1) + " KB/s"
        return (v / 1048576).toFixed(1) + " MB/s"
    }

    function getStatusText() {
        if (connStatus === "connecting") return "Connecting..."
        if (connStatus === "connected") return "Connected"
        if (connStatus === "disconnecting") return "Disconnecting..."
        return "Disconnected"
    }

    function getStatusColor() {
        if (connStatus === "connected") return successColor
        if (connStatus === "connecting" || connStatus === "disconnecting") return warningColor
        return accentColor
    }

    function currentServerModel() {
        if (!vpnController) return []
        if (currentTab === 0) return vpnController.favorites
        if (currentTab === 1) return vpnController.servers
        return vpnController.custom
    }

    function selectedServerLabel() {
        if (!vpnController) return "Best Server (Auto Connect)"
        if (selectedServerID === "" || selectedServerID === "best") return "Best Server (Auto Connect)"

        var model = currentServerModel()
        for (var i = 0; i < model.length; ++i) {
            var srv = model[i]
            var key = serverKey(srv, i)
            if (key === selectedServerID) {
                var proto = serverProtocol(srv)
                return serverName(srv) + " • " + proto
            }
        }
        return "Best Server (Auto Connect)"
    }

    Timer {
        id: autoProgressTimer
        interval: 120
        repeat: true
        running: autoScanActive && connStatus === "connecting"
        onTriggered: {
            if (autoScanProgress < 92) autoScanProgress += 2
        }
    }

    Connections {
        target: (typeof trayIcon !== "undefined") ? trayIcon : null
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
        function onStatusChanged() {
            if (connStatus === "connected") {
                autoScanActive = false
                autoScanProgress = 100
            } else if (connStatus === "disconnected") {
                autoScanActive = false
                autoScanProgress = 0
            } else if (connStatus === "disconnecting") {
                autoScanActive = false
            }
        }
        function onTrafficChanged() { }
        function onServersChanged() { }
        function onIpChanged() { }
        function onTimeChanged() { }
    }

    Menu {
        id: serverMenu
        background: Rectangle {
            color: surfaceColor
            border.color: borderColor
            radius: 14
        }
        MenuItem { text: "Refresh Servers"; onTriggered: vpnController.fetchServers() }
        MenuItem { text: "Update Subscription"; onTriggered: { vpnController.addSubscription(subscriptionUrl); vpnController.fetchServers() } }
        MenuItem { text: "Test Ping"; onTriggered: vpnController.testLatency() }
        MenuItem { text: "Copy Logs"; onTriggered: vpnController.copyLogs() }
        MenuItem { text: "Clear Servers"; onTriggered: vpnController.clearServers() }
        MenuItem { text: "Quit"; onTriggered: Qt.quit() }
    }

    Rectangle {
        anchors.fill: parent
        color: bgColor

        ColumnLayout {
            anchors.fill: parent
            anchors.margins: 24
            spacing: 18

            RowLayout {
                Layout.fillWidth: true

                Text {
                    text: "Kingo VPN"
                    color: textColor
                    font.pixelSize: 34
                    font.weight: Font.Bold
                    letterSpacing: 0.5
                }

                Item { Layout.fillWidth: true }

                Rectangle {
                    width: 36
                    height: 36
                    radius: 18
                    color: "transparent"

                    Text {
                        anchors.centerIn: parent
                        text: "⚙"
                        color: textColor
                        font.pixelSize: 22
                    }

                    MouseArea {
                        anchors.fill: parent
                        cursorShape: Qt.PointingHandCursor
                        onClicked: showSettings = !showSettings
                    }
                }
            }

            ColumnLayout {
                Layout.fillWidth: true
                spacing: 8

                Text {
                    Layout.alignment: Qt.AlignHCenter
                    text: getStatusText()
                    color: getStatusColor()
                    font.pixelSize: 23
                    font.weight: Font.Bold
                }

                Text {
                    Layout.alignment: Qt.AlignHCenter
                    text: connStatus === "connected" ? "" : "Tap to connect"
                    color: subTextColor
                    font.pixelSize: 13
                    opacity: 0.85
                }

                Item {
                    Layout.alignment: Qt.AlignHCenter
                    width: 214
                    height: 214

                    Rectangle {
                        anchors.fill: parent
                        radius: width / 2
                        color: "transparent"
                        border.color: getStatusColor()
                        border.width: 4
                        opacity: connStatus === "connecting" ? 0.8 : 1.0

                        SequentialAnimation on opacity {
                            running: connStatus === "connecting"
                            loops: Animation.Infinite
                            NumberAnimation { to: 0.25; duration: 700 }
                            NumberAnimation { to: 0.85; duration: 700 }
                        }
                    }

                    Rectangle {
                        id: innerCircle
                        anchors.centerIn: parent
                        width: 188
                        height: 188
                        radius: width / 2
                        color: surface2Color
                        border.color: borderColor
                        border.width: 1
                        scale: powerMouseArea.pressed ? 0.96 : 1.0
                        Behavior on scale { NumberAnimation { duration: 120; easing.type: Easing.OutCubic } }

                        Rectangle {
                            anchors.fill: parent
                            radius: parent.radius
                            color: getStatusColor()
                            opacity: connStatus === "connecting" ? 0.12 : 0.0
                            SequentialAnimation on opacity {
                                running: connStatus === "connecting"
                                loops: Animation.Infinite
                                NumberAnimation { to: 0.02; duration: 700 }
                                NumberAnimation { to: 0.12; duration: 700 }
                            }
                        }

                        Text {
                            anchors.centerIn: parent
                            text: "⏻"
                            color: "#FFFFFF"
                            font.pixelSize: 64
                        }

                        MouseArea {
                            id: powerMouseArea
                            anchors.fill: parent
                            cursorShape: Qt.PointingHandCursor
                            onClicked: {
                                if (connStatus === "connected") {
                                    vpnController.disconnectVpn()
                                } else if (connStatus !== "connecting" && connStatus !== "disconnecting") {
                                    selectedServerID = "best"
                                    autoScanActive = true
                                    autoScanProgress = 0
                                    vpnController.autoConnect()
                                }
                            }
                        }
                    }
                }

                Text {
                    Layout.alignment: Qt.AlignHCenter
                    text: selectedServerLabel()
                    color: accentColor
                    font.pixelSize: 14
                    font.weight: Font.Bold
                    MouseArea {
                        anchors.fill: parent
                        cursorShape: Qt.PointingHandCursor
                        onClicked: {
                            selectedServerID = "best"
                            autoScanActive = true
                            autoScanProgress = 0
                            vpnController.autoConnect()
                        }
                    }
                }

                Rectangle {
                    id: layoutAutoProgress
                    visible: autoScanActive || connStatus === "connecting"
                    Layout.fillWidth: true
                    Layout.preferredHeight: 58
                    radius: 18
                    color: surfaceColor
                    border.color: borderColor
                    border.width: 1

                    ColumnLayout {
                        anchors.fill: parent
                        anchors.margins: 12
                        spacing: 6

                        RowLayout {
                            Layout.fillWidth: true

                            Text {
                                id: tvAutoStatus
                                text: connStatus === "connecting" ? "Testing active servers..." : "Auto connect"
                                color: textColor
                                font.pixelSize: 14
                                font.weight: Font.Bold
                            }

                            Item { Layout.fillWidth: true }

                            Text {
                                id: tvAutoProgress
                                text: autoScanProgress > 0 ? autoScanProgress + "%" : ""
                                color: accentColor
                                font.pixelSize: 14
                                font.weight: Font.Bold
                            }
                        }

                        Rectangle {
                            Layout.fillWidth: true
                            height: 10
                            radius: 5
                            color: surface2Color
                            border.color: borderColor
                            border.width: 1
                            clip: true

                            Rectangle {
                                width: parent.width * (autoScanProgress / 100.0)
                                height: parent.height
                                radius: parent.radius
                                color: accentColor
                                Behavior on width { NumberAnimation { duration: 120 } }
                            }
                        }
                    }
                }
            }

            RowLayout {
                Layout.fillWidth: true
                spacing: 12

                Rectangle {
                    Layout.fillWidth: true
                    Layout.preferredHeight: 72
                    radius: 18
                    color: surfaceColor
                    border.color: borderColor
                    border.width: 1

                    ColumnLayout {
                        anchors.fill: parent
                        anchors.margins: 14
                        spacing: 4

                        Text {
                            text: "Connection Time"
                            color: subTextColor
                            font.pixelSize: 12
                        }
                        Text {
                            id: tvTimer
                            text: vpnController ? vpnController.connectionTime : "00:00:00"
                            color: textColor
                            font.pixelSize: 18
                            font.weight: Font.Bold
                        }
                    }
                }

                Rectangle {
                    Layout.fillWidth: true
                    Layout.preferredHeight: 72
                    radius: 18
                    color: surfaceColor
                    border.color: borderColor
                    border.width: 1

                    ColumnLayout {
                        anchors.fill: parent
                        anchors.margins: 14
                        spacing: 4

                        Text {
                            text: "VPN IP"
                            color: subTextColor
                            font.pixelSize: 12
                        }
                        Text {
                            id: tvIP
                            text: vpnController ? vpnController.ipAddress : "0.0.0.0"
                            color: textColor
                            font.pixelSize: 18
                            font.weight: Font.Bold
                            elide: Text.ElideRight
                        }
                    }
                }
            }

            RowLayout {
                Layout.fillWidth: true
                spacing: 12

                Rectangle {
                    Layout.fillWidth: true
                    Layout.preferredHeight: 68
                    radius: 18
                    color: surfaceColor
                    border.color: borderColor
                    border.width: 1

                    ColumnLayout {
                        anchors.fill: parent
                        anchors.margins: 14
                        spacing: 4

                        Text {
                            text: "Download"
                            color: subTextColor
                            font.pixelSize: 12
                        }
                        Text {
                            id: tvDownload
                            text: vpnController ? formatSpeed(vpnController.downloadSpeed) : "0.0 B/s"
                            color: successColor
                            font.pixelSize: 17
                            font.weight: Font.Bold
                        }
                    }
                }

                Rectangle {
                    Layout.fillWidth: true
                    Layout.preferredHeight: 68
                    radius: 18
                    color: surfaceColor
                    border.color: borderColor
                    border.width: 1

                    ColumnLayout {
                        anchors.fill: parent
                        anchors.margins: 14
                        spacing: 4

                        Text {
                            text: "Upload"
                            color: subTextColor
                            font.pixelSize: 12
                        }
                        Text {
                            id: tvUpload
                            text: vpnController ? formatSpeed(vpnController.uploadSpeed) : "0.0 B/s"
                            color: warningColor
                            font.pixelSize: 17
                            font.weight: Font.Bold
                        }
                    }
                }
            }

            Rectangle {
                Layout.fillWidth: true
                Layout.fillHeight: true
                radius: 22
                color: surfaceColor
                border.color: borderColor
                border.width: 1

                ColumnLayout {
                    anchors.fill: parent
                    anchors.margins: 18
                    spacing: 14

                    RowLayout {
                        Layout.fillWidth: true

                        Text {
                            text: "Server List"
                            color: textColor
                            font.pixelSize: 28
                            font.weight: Font.Bold
                        }

                        Item { Layout.fillWidth: true }

                        Rectangle {
                            width: 34
                            height: 34
                            radius: 17
                            color: "transparent"

                            Text {
                                anchors.centerIn: parent
                                text: "⋮"
                                color: subTextColor
                                font.pixelSize: 24
                            }

                            MouseArea {
                                anchors.fill: parent
                                cursorShape: Qt.PointingHandCursor
                                onClicked: serverMenu.popup()
                            }
                        }
                    }

                    Rectangle {
                        Layout.fillWidth: true
                        height: 46
                        radius: 16
                        color: surface2Color
                        border.color: borderColor
                        border.width: 1

                        RowLayout {
                            anchors.fill: parent
                            anchors.margins: 4
                            spacing: 4

                            Repeater {
                                model: ["Favorites", "Servers", "Custom"]
                                delegate: Rectangle {
                                    Layout.fillWidth: true
                                    Layout.fillHeight: true
                                    radius: 12
                                    color: currentTab === index ? surfaceColor : "transparent"
                                    border.color: currentTab === index ? accentColor : "transparent"
                                    border.width: 1

                                    Text {
                                        anchors.centerIn: parent
                                        text: modelData
                                        color: currentTab === index ? textColor : subTextColor
                                        font.pixelSize: 12
                                        font.weight: currentTab === index ? Font.Bold : Font.Normal
                                    }

                                    MouseArea {
                                        anchors.fill: parent
                                        cursorShape: Qt.PointingHandCursor
                                        onClicked: {
                                            currentTab = index
                                            if (currentTab !== 1) autoScanActive = false
                                        }
                                    }
                                }
                            }
                        }
                    }

                    Rectangle {
                        Layout.fillWidth: true
                        Layout.preferredHeight: 48
                        radius: 18
                        gradient: Gradient {
                            orientation: Gradient.Horizontal
                            GradientStop { position: 0.0; color: accentColor }
                            GradientStop { position: 1.0; color: accent2Color }
                        }
                        scale: getServersMouseArea.pressed ? 0.985 : 1.0
                        Behavior on scale { NumberAnimation { duration: 100 } }

                        Text {
                            anchors.centerIn: parent
                            text: "Get Active Servers"
                            color: "#FFFFFF"
                            font.pixelSize: 14
                            font.weight: Font.Bold
                        }

                        MouseArea {
                            id: getServersMouseArea
                            anchors.fill: parent
                            cursorShape: Qt.PointingHandCursor
                            onClicked: {
                                vpnController.addSubscription(subscriptionUrl)
                                vpnController.fetchServers()
                                vpnController.testLatency()
                            }
                        }
                    }

                    ListView {
                        id: serverListView
                        Layout.fillWidth: true
                        Layout.fillHeight: true
                        clip: true
                        spacing: 10
                        model: currentTab === 0 ? (vpnController ? vpnController.favorites : [])
                              : currentTab === 1 ? (vpnController ? vpnController.servers : [])
                              : (vpnController ? vpnController.custom : [])

                        ScrollBar.vertical: ScrollBar { policy: ScrollBar.AsNeeded }

                        delegate: Rectangle {
                            width: serverListView.width
                            height: 66
                            radius: 16
                            color: selectedServerID === serverKey(modelData, index) ? "#243245" : surface2Color
                            border.color: selectedServerID === serverKey(modelData, index) ? accentColor : "transparent"
                            border.width: 1
                            scale: serverMouseArea.pressed ? 0.985 : 1.0
                            Behavior on scale { NumberAnimation { duration: 100 } }
                            Behavior on color { ColorAnimation { duration: 160 } }

                            RowLayout {
                                anchors.fill: parent
                                anchors.margins: 14
                                spacing: 10

                                Rectangle {
                                    width: 10
                                    height: 10
                                    radius: 5
                                    color: getPingColor(serverLatency(modelData))
                                    Layout.alignment: Qt.AlignVCenter
                                }

                                ColumnLayout {
                                    Layout.fillWidth: true
                                    spacing: 2

                                    Text {
                                        text: serverName(modelData)
                                        color: textColor
                                        font.pixelSize: 14
                                        font.weight: Font.Bold
                                        elide: Text.ElideRight
                                        Layout.fillWidth: true
                                    }

                                    Text {
                                        text: serverProtocol(modelData) + (serverAddress(modelData) !== "" ? " • " + serverAddress(modelData) : "")
                                        color: subTextColor
                                        font.pixelSize: 10
                                        elide: Text.ElideRight
                                        Layout.fillWidth: true
                                    }
                                }

                                Text {
                                    text: serverLatency(modelData) === 9999 ? "Timeout" : serverLatency(modelData) + "ms"
                                    color: getPingColor(serverLatency(modelData))
                                    font.pixelSize: 12
                                    font.weight: Font.Bold
                                    Layout.alignment: Qt.AlignVCenter
                                }
                            }

                            MouseArea {
                                id: serverMouseArea
                                anchors.fill: parent
                                cursorShape: Qt.PointingHandCursor
                                onClicked: {
                                    selectedServerID = serverKey(modelData, index)
                                    autoScanActive = false
                                    if (modelData && modelData.uri) {
                                        vpnController.connectToServer(modelData.uri)
                                    }
                                }
                            }
                        }
                    }
                }
            }
        }

        Rectangle {
            id: settingsView
            anchors.fill: parent
            visible: showSettings
            color: "#F4F6FA"
            z: 100

            Rectangle {
                anchors.fill: parent
                color: "#000000"
                opacity: 0.22
                MouseArea {
                    anchors.fill: parent
                    onClicked: showSettings = false
                }
            }

            Rectangle {
                width: Math.min(520, parent.width - 32)
                height: Math.min(640, parent.height - 32)
                anchors.centerIn: parent
                radius: 24
                color: "#FFFFFF"
                border.color: "#E5EAF2"
                border.width: 1

                ColumnLayout {
                    anchors.fill: parent
                    anchors.margins: 20
                    spacing: 18

                    RowLayout {
                        Layout.fillWidth: true

                        Text {
                            text: "Settings"
                            color: "#111827"
                            font.pixelSize: 28
                            font.weight: Font.Bold
                        }

                        Item { Layout.fillWidth: true }

                        Rectangle {
                            width: 34
                            height: 34
                            radius: 17
                            color: "transparent"

                            Text {
                                anchors.centerIn: parent
                                text: "✕"
                                color: "#6B7280"
                                font.pixelSize: 18
                            }

                            MouseArea {
                                anchors.fill: parent
                                cursorShape: Qt.PointingHandCursor
                                onClicked: showSettings = false
                            }
                        }
                    }

                    Rectangle { Layout.fillWidth: true; height: 1; color: "#E5EAF2" }

                    ColumnLayout {
                        Layout.fillWidth: true
                        spacing: 8

                        Text {
                            text: "Split Tunneling"
                            color: "#111827"
                            font.pixelSize: 20
                            font.weight: Font.Bold
                        }

                        RowLayout {
                            Layout.fillWidth: true
                            Text {
                                text: "Choose which apps use the VPN"
                                color: "#6B7280"
                                font.pixelSize: 13
                                Layout.fillWidth: true
                                wrapMode: Text.WordWrap
                            }

                            Switch { checked: false }
                        }
                    }

                    Rectangle { Layout.fillWidth: true; height: 1; color: "#E5EAF2" }

                    ColumnLayout {
                        Layout.fillWidth: true
                        spacing: 8

                        Text {
                            text: "About"
                            color: "#111827"
                            font.pixelSize: 20
                            font.weight: Font.Bold
                        }

                        Text {
                            text: "Kingo VPN v0.5"
                            color: "#6B7280"
                            font.pixelSize: 13
                        }

                        Text {
                            text: "GitHub (Project source code)"
                            color: accentColor
                            font.pixelSize: 13
                            MouseArea {
                                anchors.fill: parent
                                cursorShape: Qt.PointingHandCursor
                            }
                        }

                        Text {
                            text: "Telegram (Join our channel)"
                            color: accentColor
                            font.pixelSize: 13
                            MouseArea {
                                anchors.fill: parent
                                cursorShape: Qt.PointingHandCursor
                            }
                        }
                    }

                    Item { Layout.fillHeight: true }

                    Text {
                        text: "by Kingo team"
                        color: "#9CA3AF"
                        font.pixelSize: 12
                        Layout.alignment: Qt.AlignHCenter
                    }
                }
            }
        }
    }

    onClosing: (close) => {
        close.accepted = false
        hide()
        if (typeof trayIcon !== "undefined" && trayIcon && trayIcon.showMessage) {
            trayIcon.showMessage("Kingo VPN", "Minimized to tray.")
        }
    }
}
