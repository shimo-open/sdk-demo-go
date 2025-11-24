import { useEffect, useState } from "react";
import { useModel } from "umi";
import { Modal, List, Popconfirm, Button } from "antd";
import { UsergroupAddOutlined } from '@ant-design/icons';
import styles from '../Team.module.less'
import { teamService } from '@/services/team.service'

export function JoinTeamModal({ open, onClose }: { open: boolean, onClose: () => void }) {
  const { joinTeam } = useModel('team');

  const [allTeams, setAllTeams] = useState([])

  useEffect(() => {
    console.log('JoinTeamModal loadAllTeams ')
    teamService.findAllTeam().then(res => {
      setAllTeams(res.data)
    }).catch(() => {
      setAllTeams([])
    })
    return () => {
      console.log('JoinTeamModal resetAllTeams ')
      setAllTeams([])
    }
  }, [])

  function confirmJoinTeam(teamId: string) {
    joinTeam(teamId)
    onClose()
  }

  return (
    <Modal
      open={open}
      title="加入团队"
      okText=""
      footer={<Button onClick={onClose}>关闭</Button>}
      onCancel={onClose}
      width="80vh"
      className={styles.modal}
    >
      <List>
        {allTeams.map((t: any, i: number) => {
          return (
            <List.Item key={i}>
              <List.Item.Meta
                avatar={<UsergroupAddOutlined />}
                title={t.name}
              />
              <Popconfirm
                title="确认加入"
                description={`确认加入团队 「${t.name}」 到团队 ?`}
                onConfirm={() => confirmJoinTeam(t.id)}
                okText="确认"
                cancelText="再想想"
              >
                <Button>加入</Button>
              </Popconfirm>
            </List.Item>
          )
        })}
      </List>
    </Modal>
  )
}