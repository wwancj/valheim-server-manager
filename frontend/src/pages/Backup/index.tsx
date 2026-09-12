import {Card, Typography} from 'antd'

const {Title} = Typography

function Backup() {
    return (
        <div>
            <Title level={2}>备份管理</Title>
            <Card>
                <p>暂无备份记录。服务器配置完成后可创建备份。</p>
            </Card>
        </div>
    )
}

export default Backup
