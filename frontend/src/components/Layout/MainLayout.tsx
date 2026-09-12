import {Layout} from 'antd'
import {Outlet} from 'react-router-dom'
import Sidebar from './Sidebar'

const {Content} = Layout

function MainLayout() {
    return (
        <Layout style={{height: '100vh', overflow: 'hidden'}}>
            <Sidebar/>
            <Layout>
                <Content style={{margin: 16, padding: 24, overflow: 'auto', height: '100vh'}}>
                    <Outlet/>
                </Content>
            </Layout>
        </Layout>
    )
}

export default MainLayout
