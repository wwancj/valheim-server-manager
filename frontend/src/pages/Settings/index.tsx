import {Card, Typography} from 'antd'

const {Title} = Typography

function Settings() {
    return (
        <div>
            <Title level={2}>设置</Title>
            <Card>
                <p>应用设置（语言、主题、数据目录等）将在后续版本中开放。</p>
            </Card>
        </div>
    )
}

export default Settings
