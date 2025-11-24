import { useModel } from 'umi'
import { useEffect, useState } from "react";
import { DataType } from "@/services/team.service";
import { Tabs, Card, Divider, Tag, Space, Button, Popconfirm } from "antd";
import { UserAddOutlined, SwapOutlined, LogoutOutlined, GoldFilled, UserOutlined } from '@ant-design/icons';
import type { TabsProps } from 'antd';
import styles from './Team.module.less'
import { AddTeamMemberModal } from './component/AddTeamMemberModal'
import { TransferTeamModal } from './component/TransferTeamModal'
import { DepartmentInfo } from './component/DepartmentInfo'

export function TeamInfo({ team, currentUserId }: { team: DataType, currentUserId: string }) {
  const { members, leaveTeam, membersDirty, teams, getMembersByTeam } = useModel('team');

  const [showAddMemberModal, setShowAddMemberModal] = useState(false)
  const [showTransferTeamModal, setShowTransferTeamModal] = useState(false)
  useEffect(() => {
    if (membersDirty && teams && teams.length > 0) {
      getMembersByTeam(teams[0].id)
    }
  }, [membersDirty])

  // Leave the team
  function doLeaveTeam(teamId: string) {
    leaveTeam(teamId)
  }

  const items: TabsProps['items'] = [
    {
      key: 'teamInfo',
      label: '团队信息',
      icon: <UserOutlined />,
      children:
        <Card>
          <h3>团队名称</h3>
          <span>{team.name}</span>
          <Divider></Divider>
          <h3>创建者</h3>
          <Tag>
            <img
              alt={`${team.creator?.name}'s avatar`}
              className={styles.avatar}
              src={team.creator?.avatar}
            />
            {team.creator.name}
            {Number(currentUserId) === Number(team.creator.id)
              ? ' (当前用户)'
              : ''}
          </Tag>
          <Divider></Divider>
          <h3>操作</h3>
          <Space>
            {Number(currentUserId) === Number(team.creator.id) ? (
              <>
                <Button icon={<UserAddOutlined />} type="primary" onClick={() => setShowAddMemberModal(true)}>添加成员</Button>
                <Button icon={<SwapOutlined />} danger onClick={() => setShowTransferTeamModal(true)}>移交团队</Button>
              </>
            ) : (
              <Popconfirm
                title="退出团队"
                description="确认离开团队？"
                onConfirm={() => doLeaveTeam(team.id)}
                okText="确认"
                cancelText="再想想"
              >
                <Button icon={<LogoutOutlined />} danger>退出团队</Button>
              </Popconfirm>
            )}
          </Space>
          <Divider></Divider>
          <h3>成员列表</h3>
          <Space>
            {members.map((member: any, index: number) => {
              return (
                <Tag key={index}>
                  <img
                    alt={`${member?.name}'s avatar`}
                    className={styles.avatar}
                    src={member?.avatar}
                  />
                  {member.name}
                  {Number(currentUserId) === Number(member.id)
                    ? ' (当前用户)'
                    : ''}
                </Tag>
              )
            })}
          </Space>
          <Divider></Divider>
          <AddTeamMemberModal open={showAddMemberModal} onClose={() => setShowAddMemberModal(false)} />
          <TransferTeamModal open={showTransferTeamModal} onClose={() => setShowTransferTeamModal(false)} />
        </Card>
    },
    {
      key: 'deptInfo',
      label: '部门信息',
      children: <DepartmentInfo />,
      icon: <GoldFilled />
    }
  ]
  return (
    <div className={styles.tabBody}>
      <Tabs defaultActiveKey="teamInfo" items={items} />
    </div>
  )
}