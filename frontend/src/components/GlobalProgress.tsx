import {Progress, Card, Typography} from 'antd'
import {CheckCircleOutlined, LoadingOutlined} from '@ant-design/icons'
import {useProgress} from '../contexts/ProgressContext'

const {Text} = Typography

export default function GlobalProgress() {
    const {progress} = useProgress()

    if (!progress.visible) return null

    return (
        <Card
            size="small"
            style={{
                position: 'fixed',
                bottom: 20,
                right: 20,
                width: 360,
                zIndex: 1000,
                boxShadow: '0 4px 12px rgba(0,0,0,0.4)',
                border: progress.done ? '1px solid #52c41a' : '1px solid #d4a017',
            }}
        >
            <div style={{display: 'flex', alignItems: 'center', gap: 8, marginBottom: 8}}>
                {progress.done
                    ? <CheckCircleOutlined style={{color: '#52c41a', fontSize: 16}}/>
                    : <LoadingOutlined style={{color: '#d4a017', fontSize: 16}}/>
                }
                <Text strong>{progress.title}</Text>
                {progress.done && <Text type="success" style={{marginLeft: 'auto'}}>完成</Text>}
            </div>
            <Progress
                percent={progress.percent}
                status={progress.done ? 'success' : 'active'}
                strokeColor={progress.done ? '#52c41a' : '#d4a017'}
                size="small"
            />
            <Text type="secondary" style={{fontSize: 12, marginTop: 4, display: 'block'}}>
                {progress.message}
            </Text>
        </Card>
    )
}
