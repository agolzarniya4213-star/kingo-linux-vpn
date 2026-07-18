/****************************************************************************
** Meta object code from reading C++ file 'appcore.h'
**
** Created by: The Qt Meta Object Compiler version 67 (Qt 5.15.19)
**
** WARNING! All changes made in this file will be lost!
*****************************************************************************/

#include <memory>
#include "../../../appcore.h"
#include <QtCore/qbytearray.h>
#include <QtCore/qmetatype.h>
#if !defined(Q_MOC_OUTPUT_REVISION)
#error "The header file 'appcore.h' doesn't include <QObject>."
#elif Q_MOC_OUTPUT_REVISION != 67
#error "This file was generated using the moc from 5.15.19. It"
#error "cannot be used with the include files from this version of Qt."
#error "(The moc has changed too much.)"
#endif

QT_BEGIN_MOC_NAMESPACE
QT_WARNING_PUSH
QT_WARNING_DISABLE_DEPRECATED
struct qt_meta_stringdata_AppCore_t {
    QByteArrayData data[18];
    char stringdata0[252];
};
#define QT_MOC_LITERAL(idx, ofs, len) \
    Q_STATIC_BYTE_ARRAY_DATA_HEADER_INITIALIZER_WITH_OFFSET(len, \
    qptrdiff(offsetof(qt_meta_stringdata_AppCore_t, stringdata0) + ofs \
        - idx * sizeof(QByteArrayData)) \
    )
static const qt_meta_stringdata_AppCore_t qt_meta_stringdata_AppCore = {
    {
QT_MOC_LITERAL(0, 0, 7), // "AppCore"
QT_MOC_LITERAL(1, 8, 22), // "connectionStateChanged"
QT_MOC_LITERAL(2, 31, 0), // ""
QT_MOC_LITERAL(3, 32, 21), // "selectedServerChanged"
QT_MOC_LITERAL(4, 54, 23), // "showServerListRequested"
QT_MOC_LITERAL(5, 78, 12), // "statsUpdated"
QT_MOC_LITERAL(6, 91, 16), // "onConnectSuccess"
QT_MOC_LITERAL(7, 108, 19), // "onDisconnectSuccess"
QT_MOC_LITERAL(8, 128, 14), // "onCommandError"
QT_MOC_LITERAL(9, 143, 15), // "updateMockStats"
QT_MOC_LITERAL(10, 159, 16), // "toggleConnection"
QT_MOC_LITERAL(11, 176, 14), // "openServerList"
QT_MOC_LITERAL(12, 191, 11), // "isConnected"
QT_MOC_LITERAL(13, 203, 10), // "statusText"
QT_MOC_LITERAL(14, 214, 14), // "selectedServer"
QT_MOC_LITERAL(15, 229, 6), // "rxData"
QT_MOC_LITERAL(16, 236, 6), // "txData"
QT_MOC_LITERAL(17, 243, 8) // "pingData"

    },
    "AppCore\0connectionStateChanged\0\0"
    "selectedServerChanged\0showServerListRequested\0"
    "statsUpdated\0onConnectSuccess\0"
    "onDisconnectSuccess\0onCommandError\0"
    "updateMockStats\0toggleConnection\0"
    "openServerList\0isConnected\0statusText\0"
    "selectedServer\0rxData\0txData\0pingData"
};
#undef QT_MOC_LITERAL

static const uint qt_meta_data_AppCore[] = {

 // content:
       8,       // revision
       0,       // classname
       0,    0, // classinfo
      10,   14, // methods
       6,   80, // properties
       0,    0, // enums/sets
       0,    0, // constructors
       0,       // flags
       4,       // signalCount

 // signals: name, argc, parameters, tag, flags
       1,    0,   64,    2, 0x06 /* Public */,
       3,    0,   65,    2, 0x06 /* Public */,
       4,    0,   66,    2, 0x06 /* Public */,
       5,    0,   67,    2, 0x06 /* Public */,

 // slots: name, argc, parameters, tag, flags
       6,    1,   68,    2, 0x08 /* Private */,
       7,    1,   71,    2, 0x08 /* Private */,
       8,    1,   74,    2, 0x08 /* Private */,
       9,    0,   77,    2, 0x08 /* Private */,

 // methods: name, argc, parameters, tag, flags
      10,    0,   78,    2, 0x02 /* Public */,
      11,    0,   79,    2, 0x02 /* Public */,

 // signals: parameters
    QMetaType::Void,
    QMetaType::Void,
    QMetaType::Void,
    QMetaType::Void,

 // slots: parameters
    QMetaType::Void, QMetaType::QString,    2,
    QMetaType::Void, QMetaType::QString,    2,
    QMetaType::Void, QMetaType::QString,    2,
    QMetaType::Void,

 // methods: parameters
    QMetaType::Void,
    QMetaType::Void,

 // properties: name, type, flags
      12, QMetaType::Bool, 0x00495001,
      13, QMetaType::QString, 0x00495001,
      14, QMetaType::QString, 0x00495103,
      15, QMetaType::QString, 0x00495001,
      16, QMetaType::QString, 0x00495001,
      17, QMetaType::QString, 0x00495001,

 // properties: notify_signal_id
       0,
       0,
       1,
       3,
       3,
       3,

       0        // eod
};

void AppCore::qt_static_metacall(QObject *_o, QMetaObject::Call _c, int _id, void **_a)
{
    if (_c == QMetaObject::InvokeMetaMethod) {
        auto *_t = static_cast<AppCore *>(_o);
        (void)_t;
        switch (_id) {
        case 0: _t->connectionStateChanged(); break;
        case 1: _t->selectedServerChanged(); break;
        case 2: _t->showServerListRequested(); break;
        case 3: _t->statsUpdated(); break;
        case 4: _t->onConnectSuccess((*reinterpret_cast< const QString(*)>(_a[1]))); break;
        case 5: _t->onDisconnectSuccess((*reinterpret_cast< const QString(*)>(_a[1]))); break;
        case 6: _t->onCommandError((*reinterpret_cast< const QString(*)>(_a[1]))); break;
        case 7: _t->updateMockStats(); break;
        case 8: _t->toggleConnection(); break;
        case 9: _t->openServerList(); break;
        default: ;
        }
    } else if (_c == QMetaObject::IndexOfMethod) {
        int *result = reinterpret_cast<int *>(_a[0]);
        {
            using _t = void (AppCore::*)();
            if (*reinterpret_cast<_t *>(_a[1]) == static_cast<_t>(&AppCore::connectionStateChanged)) {
                *result = 0;
                return;
            }
        }
        {
            using _t = void (AppCore::*)();
            if (*reinterpret_cast<_t *>(_a[1]) == static_cast<_t>(&AppCore::selectedServerChanged)) {
                *result = 1;
                return;
            }
        }
        {
            using _t = void (AppCore::*)();
            if (*reinterpret_cast<_t *>(_a[1]) == static_cast<_t>(&AppCore::showServerListRequested)) {
                *result = 2;
                return;
            }
        }
        {
            using _t = void (AppCore::*)();
            if (*reinterpret_cast<_t *>(_a[1]) == static_cast<_t>(&AppCore::statsUpdated)) {
                *result = 3;
                return;
            }
        }
    }
#ifndef QT_NO_PROPERTIES
    else if (_c == QMetaObject::ReadProperty) {
        auto *_t = static_cast<AppCore *>(_o);
        (void)_t;
        void *_v = _a[0];
        switch (_id) {
        case 0: *reinterpret_cast< bool*>(_v) = _t->isConnected(); break;
        case 1: *reinterpret_cast< QString*>(_v) = _t->statusText(); break;
        case 2: *reinterpret_cast< QString*>(_v) = _t->selectedServer(); break;
        case 3: *reinterpret_cast< QString*>(_v) = _t->rxData(); break;
        case 4: *reinterpret_cast< QString*>(_v) = _t->txData(); break;
        case 5: *reinterpret_cast< QString*>(_v) = _t->pingData(); break;
        default: break;
        }
    } else if (_c == QMetaObject::WriteProperty) {
        auto *_t = static_cast<AppCore *>(_o);
        (void)_t;
        void *_v = _a[0];
        switch (_id) {
        case 2: _t->setSelectedServer(*reinterpret_cast< QString*>(_v)); break;
        default: break;
        }
    } else if (_c == QMetaObject::ResetProperty) {
    }
#endif // QT_NO_PROPERTIES
}

QT_INIT_METAOBJECT const QMetaObject AppCore::staticMetaObject = { {
    QMetaObject::SuperData::link<QObject::staticMetaObject>(),
    qt_meta_stringdata_AppCore.data,
    qt_meta_data_AppCore,
    qt_static_metacall,
    nullptr,
    nullptr
} };


const QMetaObject *AppCore::metaObject() const
{
    return QObject::d_ptr->metaObject ? QObject::d_ptr->dynamicMetaObject() : &staticMetaObject;
}

void *AppCore::qt_metacast(const char *_clname)
{
    if (!_clname) return nullptr;
    if (!strcmp(_clname, qt_meta_stringdata_AppCore.stringdata0))
        return static_cast<void*>(this);
    return QObject::qt_metacast(_clname);
}

int AppCore::qt_metacall(QMetaObject::Call _c, int _id, void **_a)
{
    _id = QObject::qt_metacall(_c, _id, _a);
    if (_id < 0)
        return _id;
    if (_c == QMetaObject::InvokeMetaMethod) {
        if (_id < 10)
            qt_static_metacall(this, _c, _id, _a);
        _id -= 10;
    } else if (_c == QMetaObject::RegisterMethodArgumentMetaType) {
        if (_id < 10)
            *reinterpret_cast<int*>(_a[0]) = -1;
        _id -= 10;
    }
#ifndef QT_NO_PROPERTIES
    else if (_c == QMetaObject::ReadProperty || _c == QMetaObject::WriteProperty
            || _c == QMetaObject::ResetProperty || _c == QMetaObject::RegisterPropertyMetaType) {
        qt_static_metacall(this, _c, _id, _a);
        _id -= 6;
    } else if (_c == QMetaObject::QueryPropertyDesignable) {
        _id -= 6;
    } else if (_c == QMetaObject::QueryPropertyScriptable) {
        _id -= 6;
    } else if (_c == QMetaObject::QueryPropertyStored) {
        _id -= 6;
    } else if (_c == QMetaObject::QueryPropertyEditable) {
        _id -= 6;
    } else if (_c == QMetaObject::QueryPropertyUser) {
        _id -= 6;
    }
#endif // QT_NO_PROPERTIES
    return _id;
}

// SIGNAL 0
void AppCore::connectionStateChanged()
{
    QMetaObject::activate(this, &staticMetaObject, 0, nullptr);
}

// SIGNAL 1
void AppCore::selectedServerChanged()
{
    QMetaObject::activate(this, &staticMetaObject, 1, nullptr);
}

// SIGNAL 2
void AppCore::showServerListRequested()
{
    QMetaObject::activate(this, &staticMetaObject, 2, nullptr);
}

// SIGNAL 3
void AppCore::statsUpdated()
{
    QMetaObject::activate(this, &staticMetaObject, 3, nullptr);
}
QT_WARNING_POP
QT_END_MOC_NAMESPACE
