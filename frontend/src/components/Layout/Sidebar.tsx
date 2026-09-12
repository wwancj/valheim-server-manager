import {Layout, Menu} from 'antd'
import {
    DashboardOutlined,
    CloudServerOutlined,
    AppstoreOutlined,
    SettingOutlined,
    FileTextOutlined,
    SaveOutlined,
    ToolOutlined,
} from '@ant-design/icons'
import {useNavigate, useLocation} from 'react-router-dom'

const {Sider} = Layout

const menuItems = [
    {key: '/dashboard', icon: <DashboardOutlined/>, label: 'Dashboard'},
    {key: '/server', icon: <CloudServerOutlined/>, label: '服务器'},
    {key: '/mods', icon: <AppstoreOutlined/>, label: 'Mods'},
    {key: '/config', icon: <ToolOutlined/>, label: '配置'},
    {key: '/logs', icon: <FileTextOutlined/>, label: '日志'},
    {key: '/backup', icon: <SaveOutlined/>, label: '备份'},
    {key: '/settings', icon: <SettingOutlined/>, label: '设置'},
]

function Sidebar() {
    const navigate = useNavigate()
    const location = useLocation()

    return (
        <Sider width={200} theme="dark">
            <div style={{
                height: 64,
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                color: '#d4a017',
                fontSize: 18,
                fontWeight: 'bold',
            }}>
                ⚔ Valheim Server Manager
            </div>
            <Menu
                theme="dark"
                mode="inline"
                selectedKeys={[location.pathname]}
                items={menuItems}
                onClick={({key}) => navigate(key)}
            />
        </Sider>
    )
}

export default Sidebar
