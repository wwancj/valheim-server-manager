import React from 'react'
import {createRoot} from 'react-dom/client'
import {ConfigProvider, theme} from 'antd'
import zhCN from 'antd/locale/zh_CN'
import App from './App'
import {wsClient} from './api/websocket'
import './style.css'

// Connect WebSocket for real-time events
wsClient.connect()

const container = document.getElementById('root')
const root = createRoot(container!)

root.render(
    <React.StrictMode>
        <ConfigProvider
            locale={zhCN}
            theme={{
                algorithm: theme.darkAlgorithm,
                token: {
                    colorPrimary: '#d4a017',
                    borderRadius: 6,
                },
            }}
        >
            <App/>
        </ConfigProvider>
    </React.StrictMode>
)
