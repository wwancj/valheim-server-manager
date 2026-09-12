import {HashRouter, Routes, Route, Navigate} from 'react-router-dom'
import MainLayout from './components/Layout/MainLayout'
import Dashboard from './pages/Dashboard'
import Server from './pages/Server'
import Mods from './pages/Mods'
import Config from './pages/Config'
import Logs from './pages/Logs'
import Backup from './pages/Backup'
import Settings from './pages/Settings'

function App() {
    return (
        <HashRouter>
            <Routes>
                <Route path="/" element={<MainLayout/>}>
                    <Route index element={<Navigate to="/dashboard" replace/>}/>
                    <Route path="dashboard" element={<Dashboard/>}/>
                    <Route path="server" element={<Server/>}/>
                    <Route path="mods" element={<Mods/>}/>
                    <Route path="config" element={<Config/>}/>
                    <Route path="logs" element={<Logs/>}/>
                    <Route path="backup" element={<Backup/>}/>
                    <Route path="settings" element={<Settings/>}/>
                </Route>
            </Routes>
        </HashRouter>
    )
}

export default App
