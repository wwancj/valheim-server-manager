import {useState, useEffect} from 'react'
import {Card, Col, Row, Statistic, Typography} from 'antd'
import {CloudServerOutlined, AppstoreOutlined, SaveOutlined} from '@ant-design/icons'
import {GetAppInfo} from '../../../wailsjs/go/app/App'

const {Title} = Typography

function Dashboard() {
    const [appInfo, setAppInfo] = useState<Record<string, string>>({})

    useEffect(() => {
        GetAppInfo().then(setAppInfo)
    }, [])

    return (
        <div>
            <Title level={2}>Dashboard</Title>
            <Row gutter={[16, 16]}>
                <Col span={8}>
                    <Card>
                        <Statistic
                            title="服务器状态"
                            value="未运行"
                            prefix={<CloudServerOutlined/>}
                        />
                    </Card>
                </Col>
                <Col span={8}>
                    <Card>
                        <Statistic
                            title="已安装 Mod"
                            value={0}
                            prefix={<AppstoreOutlined/>}
                        />
                    </Card>
                </Col>
                <Col span={8}>
                    <Card>
                        <Statistic
                            title="备份数量"
                            value={0}
                            prefix={<SaveOutlined/>}
                        />
                    </Card>
                </Col>
            </Row>
            <Card style={{marginTop: 16}}>
                <Title level={4}>应用信息</Title>
                <p>名称：{appInfo.name || '-'}</p>
                <p>版本：{appInfo.version || '-'}</p>
                <p>Go：{appInfo.go || '-'}</p>
                <p>系统：{appInfo.os || '-'} / {appInfo.arch || '-'}</p>
            </Card>
        </div>
    )
}

export default Dashboard
