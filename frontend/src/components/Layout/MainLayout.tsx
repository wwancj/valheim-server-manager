import {Layout} from 'antd'
import {Outlet} from 'react-router-dom'
import Sidebar from './Sidebar'

const {Content} = Layout

function MainLayout() {
    return (
        <Layout style={{minHeight: '100vh'}}>
            <Sidebar/>
            <Layout>
                <Content style={{margin: 16, padding: 24}}>
                    <Outlet/>
                </Content>
            </Layout>
        </Layout>
    )
}

export default MainLayout
