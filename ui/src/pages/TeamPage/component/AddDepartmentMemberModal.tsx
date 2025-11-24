import { useEffect, useState } from "react";
import { useModel } from "umi";
import { Modal, List, Popconfirm, Button, message } from "antd";
import { UserOutlined } from '@ant-design/icons';
import styles from '../Team.module.less'
import { teamService } from "@/services/team.service";

export function AddDepartmentMemberModal({ onCreated, departmentId, open, onClose }: { onCreated: () => void, departmentId: number, open: boolean, onClose: () => void }) {
  const [confirmLoading, setConfirmLoading] = useState(false)
  const [deptUsers, setDeptUsers] = useState([])
  const [deptUsersLoaded, setDeptUsersLoaded] = useState(false)
  const [lastDepartmentId, setLastDepartmentId] = useState(departmentId)

  const { userInfo } = useModel('user');
  const { teams, members } = useModel('team');

  const { currentUserId, teamMembers } = {
    currentUserId: userInfo?.id,
    teamMembers: members
  }

  const currentTeamId = teams && teams.length > 0 ? teams[0].id : null

  const deptMembersMap = deptUsers.reduce((acc: any, member: any) => {
    acc[member.id] = member
    return acc
  }, {})
  // Fetch the department members
  useEffect(() => {
    if (!deptUsersLoaded && departmentId && lastDepartmentId !== departmentId) {
      teamService.getDepartmentMembers(currentTeamId, departmentId).then((members) => {
        setDeptUsers(members.data)
        setDeptUsersLoaded(true)
        setLastDepartmentId(departmentId)
      })
    }
  }, [deptUsersLoaded, departmentId])

  useEffect(() => {
    if (deptUsersLoaded && lastDepartmentId !== departmentId) {
      setDeptUsersLoaded(false)
    }
  }, [deptUsersLoaded, lastDepartmentId, departmentId])

  // Confirm adding the user
  function confirmAction(id: string) {
    if (currentTeamId && departmentId) {
      setConfirmLoading(true)
      teamService.updateUserDepartment(currentTeamId, id, departmentId).then(res => {
        onCreated()
      }).finally(() => {
        onClose()
        setConfirmLoading(false)
      })
    } else {
      message.warning('团队或部门信息未找到')
    }
  }

  return (
    <Modal
      open={open}
      title="添加用户到部门"
      okText=""
      footer={<Button onClick={onClose}>关闭</Button>}
      onCancel={onClose}
      width="80vh"
      className={styles.modal}
    >
      <List>
        {(teamMembers || []).map((m: any, i: number) => {
          const hasJoined = !!deptMembersMap[m.id]
          const isCurrentUser = Number(m.id) === Number(currentUserId)
          return (
            <List.Item key={i}>
              <List.Item.Meta
                avatar={<UserOutlined />}
                title={<>
                  {m.name} - {m.email}
                  {hasJoined ? ' (已加入)' : ''}
                  {isCurrentUser ? ' (当前用户)' : ''}
                </>}
              />
              <Popconfirm
                title="确认添加"
                description={`确认添加用户 「${m.name}」 到此部门 ? 将会退出其他部门`}
                onConfirm={() => confirmAction(m.id)}
                okText="确认"
                cancelText="再想想"
              >
                <Button style={{ visibility: !hasJoined ? 'visible' : 'hidden' }} loading={confirmLoading}>
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