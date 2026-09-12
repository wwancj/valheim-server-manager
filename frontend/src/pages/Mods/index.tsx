import {Card, Typography, Input, Tabs} from 'antd'
import {SearchOutlined} from '@ant-design/icons'

const {Title} = Typography

function Mods() {
    return (
        <div>
            <Title level={2}>Mod 管理</Title>
            <Tabs
                items={[
                    {
                        key: 'installed',
                        label: '已安装',
                        children: (
                            <Card>
                                <p>暂无已安装的 Mod</p>
                            </Card>
                        ),
                    },
                    {
                        key: 'store',
                        label: 'Mod 商店',
                        children: (
                            <Card>
                                <Input
                                    placeholder="搜索 Mod..."
                                    prefix={<SearchOutlined/>}
                                    style={{marginBottom: 16}}
                                />
                                <p>连接 Thunderstore 后可搜索安装 Mod</p>
                            </Card>
                        ),
                    },
                ]}
            />
        </div>
    )
}

export default Mods
