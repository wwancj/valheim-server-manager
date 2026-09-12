import {Card, Typography} from 'antd'

const {Title} = Typography

function Logs() {
    return (
        <div>
            <Title level={2}>实时日志</Title>
            <Card>
                <p>服务器启动后，日志将在此实时显示。</p>
            </Card>
        </div>
    )
}

export default Logs
