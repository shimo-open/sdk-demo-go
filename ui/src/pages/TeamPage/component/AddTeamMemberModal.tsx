import { useEffect, useState } from "react";
import { useModel } from "umi";
import { Modal, List, Popconfirm, Button, message } from "antd";
import { teamContants } from "@/constants";
import { UserOutlined } from '@ant-design/icons';
import styles from '../Team.module.less';
import { userService } from '@/services/user.service'

export function AddTeamMemberModal({ open, onClose }: { open: boolean, onClose: () => void }) {
  const { teams, members, addMember } = useModel('team')
  const { userInfo } = useModel('user')

  const [allUsers, setAllUsers] = useState([])
  const { currentUserId } = {
    currentUserId: userInfo?.id
  }
  const currentTeamId = teams && teams.length > 0 ? teams[0].id : null

  const membersMap = members.reduce((acc: any, member: any) => {
    acc[member.id] = member
    return acc
  }, {})

  useEffect(() => {
    userService.getAll().then(res => {
      setAllUsers(res.data)
    }).catch(err => {
      setAllUsers([])
    })
  }, [])

  // Add a user to the team
  function confirmAction(id: string) {
    if (currentTeamId) {
      addMember({ teamId: currentTeamId, userId: id })
      onClose()
    } else {
      message.warning('未找到团队信息')
    }
  }

  return (
    <Modal
      open={open}
      title="添加用户到团队"
      okText=""
      footer={<Button onClick={onClose}>关闭</Button>}
      onCancel={onClose}
      width="80vh"
      className={styles.modal}
    >
      <List>
        {(allUsers || []).map((m: any, i: number) => {
          const hasJoined = !!membersMap[m.id]
          const isCurrentUser = Number(m.id) === Number(currentUserId)
          const joinedOtherTeam = m.teamId && m.teamId !== currentTeamId
          const canBeAdded =
            !hasJoined && !isCurrentUser && !joinedOtherTeam
          return (
            <List.Item key={i}>
              <List.Item.Meta
                avatar={<UserOutlined />}
                title={<>
                  {m.name} - {m.email}
                  {isCurrentUser
                    ? ' (当前用户)'
                    : hasJoined
                      ? ' (已加入)'
                      : joinedOtherTeam
                        ? ' (已加入其他团队)'
                        : ''}
                </>}
              />
              <Popconfirm
                title="确认添加"
                description={`确认添加用户 「${m.name}」 到团队 ?`}
                onConfirm={() => confirmAction(m.id)}
                okText="确认"
                cancelText="再想想"
              >
                <Button style={{ visibility: canBeAdded ? 'visible' : 'hidden' }}>
                  添加
                </Button>
              </Popconfirm>
            </List.Item>
          )
        })}
      </List>
    </Modal>
  )
}