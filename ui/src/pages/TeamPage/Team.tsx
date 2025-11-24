import React, { useEffect } from 'react'
import styles from './Team.module.less'
import { useModel } from "umi";
import { Card, Spin, Button, Space } from "antd";
import { PlusOutlined, UsergroupAddOutlined } from '@ant-design/icons';
import { CreateTeamModal } from './component/CreateTeamModal'
import { JoinTeamModal } from './component/JoinTeamModal'
import { TeamInfo } from './TeamInfo'

export default function Team() {
  const { userInfo } = useModel('user');
  const { teamLoaded, teamLoading, dirty, teams, membersLoading, membersLoaded, getMineTeams, resetTeams, getMembersByTeam } = useModel('team');
  const { loading, loaded, userId } = {
    loading: teamLoading,
    loaded: teamLoaded,
    userId: userInfo?.id
  }
  const [showCreateTeamModal, setShowCreateTeamModal] = React.useState(false)
  const [showJoinTeamModal, setShowJoinTeamModal] = React.useState(false)

  // Load the teams list
  useEffect(() => {
    if (!loading && !loaded && userId) {
      console.log('Team loadMineTeams')
      getMineTeams(userId)
    }
    return () => {
      console.log('Team resetTeams')
      resetTeams()
    }
  }, [userId])

  // Load the member list
  useEffect(() => {
    // Only one team is supported; use the first entry
    if (!membersLoading && !membersLoaded && teams && teams.length > 0) {
      console.log('Team loadTeamMembers')
      getMembersByTeam(teams[0].id)
    }
  }, [teams])

  // Reload when team data changes
  useEffect(() => {
    if (dirty) {
      console.log('Team loadMineTeams dirty')
      getMineTeams(userId)
    }
  }, [dirty, userId])

  return (
    <>
      <Space className={styles.actionButtonContainer}>
        <Button
          icon={<PlusOutlined />}
          disabled={teams.length > 0}
          onClick={() => setShowCreateTeamModal(true)}
        >
          创建团队
        </Button>
        <Button
          icon={<UsergroupAddOutlined />}
          disabled={teams.length > 0}
          onClick={() => setShowJoinTeamModal(true)}
        >
          加入团队
        </Button>
      </Space>
      <Spin spinning={loading}>
        {teams && teams.length > 0 ? (
          <TeamInfo team={teams[0]} currentUserId={userId} />
        ) : (
          <Card className={styles.actionBodyContainer}>
            还未加入任何团队
          </Card>
        )}
      </Spin>
      <CreateTeamModal open={showCreateTeamModal} onClose={() => setShowCreateTeamModal(false)} />
      <JoinTeamModal open={showJoinTeamModal} onClose={() => setShowJoinTeamModal(false)} />
    </>
  )
}
