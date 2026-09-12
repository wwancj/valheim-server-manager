import {Card, Typography} from 'antd'

const {Title} = Typography

function Config() {
    return (
        <div>
            <Title level={2}>服务器配置</Title>
            <Card>
                <p>请先添加并选择一个服务器，然后在此配置服务器参数。</p>
            </Card>
        </div>
    )
}

export default Config
