import {Card, Typography, Button, Table, Space} from 'antd'
import {PlusOutlined} from '@ant-design/icons'

const {Title} = Typography

function Server() {
    return (
        <div>
            <div style={{display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 16}}>
                <Title level={2} style={{margin: 0}}>服务器管理</Title>
                <Button type="primary" icon={<PlusOutlined/>}>添加服务器</Button>
            </div>
            <Card>
                <Table
                    dataSource={[]}
                    columns={[
                        {title: '名称', dataIndex: 'name', key: 'name'},
                        {title: '路径', dataIndex: 'installPath', key: 'installPath'},
                        {title: '状态', dataIndex: 'status', key: 'status'},
                        {title: '操作', key: 'action'},
                    ]}
                    locale={{emptyText: '暂无服务器，点击上方按钮添加'}}
                />
            </Card>
        </div>
    )
}

export default Server
