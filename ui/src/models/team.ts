import { teamContants } from '@/constants'
import { teamService } from "@/services/team.service";
import { message } from 'antd';
import { useState } from 'react'

const teamModel = () => {
  const [dirty, setDirty] = useState(false)
  const [teams, setTeams] = useState([] as any)

  const [teamLoading, setTeamLoading] = useState(false)
  const [teamLoaded, setTeamLoaded] = useState(false)

  const [membersLoading, setMembersLoading] = useState(false)
  const [membersLoaded, setMembersLoaded] = useState(false)
  const [members, setMembers] = useState([] as any)
  const [membersDirty, setMembersDirty] = useState(false)

  const resetTeams = () => {
    setTeamLoading(false)
    setTeamLoaded(false)
    setTeams([])
  }

  const loadedTeams = () => {
    setTeamLoading(false)
    setTeamLoaded(true)
  }

  // Fetch teams joined by the current user
  const getMineTeams = (userId: string) => {
    setTeamLoading(true)
    teamService.getMineTeams(userId).then(res => {
      setTeams(res.data)
      setDirty(false)
      setMembersDirty(true)
      loadedTeams()
    })
  }

  const loadedMembers = () => {
    setMembersLoaded(true)
    setMembersLoading(false)
  }

  // Fetch team members
  const getMembersByTeam = (teamId: string) => {
    setMembersLoading(true)
    teamService.getMembersByTeam(teamId).then(res => {
      setMembers(res.data)
      setMembersDirty(false)
      loadedMembers()
    })
  }
  // Add a team member
  const addMember = (payload: { teamId: string, userId: string }) => {
    teamService.addMember(payload).then(res => {
      setDirty(true)
    })
  }
  // Transfer team ownership
  const transferTeam = (payload: { teamId: string, userId: string }) => {
    teamService.transferTeam(payload).then(res => {
      setDirty(true)
    })
  }
  // Create a team
  const createTeam = (name: string) => {
    teamService.createTeam(name).then(res => {
      setDirty(true)
    })
  }
  // Join a team
  const joinTeam = (teamId: string) => {
    teamService.joinTeam(teamId).then(res => {
      message.success('成功加入团队');
      setDirty(true)
    })
  }

  const leaveTeam = (teamId: string) => {
    teamService.leaveTeam(teamId).then(res => {
      setDirty(true)
      message.success('离开团队成功');
    })
  }
  return {
    dirty,
    teams,
    teamLoading,
    teamLoaded,
    membersLoading,
    membersLoaded,
    members,
    membersDirty,
    resetTeams,
    getMineTeams,
    getMembersByTeam,
    addMember,
    transferTeam,
    createTeam,
    joinTeam,
    leaveTeam
  }
}
export default teamModel
