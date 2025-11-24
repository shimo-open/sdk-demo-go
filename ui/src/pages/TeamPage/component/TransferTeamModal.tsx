import { useModel } from "umi";
import { Modal, Popconfirm, List, Avatar, Button, message } from "antd";
import { fileConstants } from "@/constants";
import { UserOutlined } from '@ant-design/icons';
import styles from '../FileList.module.less'
import { teamContants } from '@/constants'

export function TransferTeamModal({ open, onClose }: { open: boolean, onClose: () => void }) {
  const { userInfo } = useModel('user');
  const { members, teams, transferTeam } = useModel('team');

  const { currentUserId } = {
    currentUserId: userInfo?.id
  }

  const currentTeam = teams && teams.length > 0 ? teams[0] : null
  const teamName = currentTeam ? currentTeam.name : '未找到团队'

  function confirmAction(id: string) {
    const teamId = currentTeam ? currentTeam.id : null
    if (teamId) {
      transferTeam({ teamId: teamId, userId: id })
      onClose()
    } else {
      message.error('未找到团队信息')
    }
  }
  return (
    <Modal
      open={open}
      title="移交团队"
      footer={<Button onClick={onClose}>关闭</Button>}
      onCancel={onClose}
    >
      <List>
        {members.map((m: any, i: number) => {
          return (
            <List.Item key={i}>
              <List.Item.Meta
                avatar={<UserOutlined />}
                title={
                  <>
                    {m.name} - {m.email}
                    {Number(currentUserId) === Number(m.id) ? ' (当前用户)' : ''}
                  </>
                }
              />
              <Popconfirm
                title="确认移交"
                description={`确认将团队 「${teamName}」 移交给用户 「${m.name}」 ?`}
                onConfirm={() => confirmAction(m.id)}
                okText="确认"
                cancelText="再想想"
              >
                <Button style={{ visibility: Number(m.id) !== Number(currentUserId) ? 'visible' : 'hidden' }}>
                  移交
                </Button>
              </Popconfirm>
            </List.Item>
          )
        })}
      </List>
    </Modal>
  )
}