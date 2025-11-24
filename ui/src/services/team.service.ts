import req from '@/utils/axios';

export const teamService = {
  createTeam,
  getMineTeams,
  findAllTeam,
  joinTeam,
  leaveTeam,
  addMember,
  transferTeam,
  getMembersByTeam,
  departmentTopTree,
  getDepartmentChildren,
  addDepartment,
  updateUserDepartment,
  removeDepartment,
  removeDepartmentMember,
  getDepartmentMembers
}

export interface DataType {
  createdAt: number;
  creatorId: number;
  creator: any;
  id: string;
  name: string;
  role: any;
  updatedAt: string;
}

// Create a team
async function createTeam(name: string) {
  return await req.post(`api/teams`, { name: name })
}

// Fetch teams the current user has joined
async function getMineTeams(userId: string | number) {
  const teams = req.get(`api/users/${userId}/teams`)
  return teams
}

// Fetch all teams
async function findAllTeam() {
  const teams = await req.get(`api/teams`)
  return teams
}

// Join a team
async function joinTeam(teamId: string | number) {
  return await req.post(`api/teams/${teamId}/members`)
}

// Leave a team
async function leaveTeam(teamId: string | number) {
  return await req.delete(`api/users/me/teams/${teamId}`)
}

// Fetch members for a team
async function getMembersByTeam(teamId: string | number) {
  const members = await req.get(`api/teams/${teamId}/members`)
  return members
}

// Transfer team ownership
async function transferTeam(transData: { teamId: string | number, userId: string | number }) {
  return await req.patch(`api/teams/${transData.teamId}/role/creator`, { newCreatorId: transData.userId })
}

// Add a team member
async function addMember(addData: { teamId: string | number, userId: string | number }) {
  return await req.post(`api/teams/${addData.teamId}/members`, { userId: addData.userId })
}

async function departmentTopTree(teamId: string | number) {
  return await req.get(`api/teams/${teamId}/department-top-tree`)
}

async function getDepartmentChildren(teamId: string | number, departmentId: string | number) {
  return await req.get(`api/teams/${teamId}/departments/${departmentId}/children`)
}

async function addDepartment({ teamId, departmentId, name }: { teamId: string | number, departmentId: string | number, name: string }) {
  return await req.post(`api/teams/${teamId}/departments`, { parentId: departmentId, name })
}

async function getDepartmentMembers(teamId: string | number, departmentId: string | number) {
  return await req.get(`api/teams/${teamId}/departments/${departmentId}/members`)
}

// Add or remove department members
async function updateUserDepartment(teamId: string | number, userId: string | number, departmentId: string | number) {
  return await req.patch(`api/teams/${teamId}/departments/${departmentId}/members`, { userId: userId })
}

async function removeDepartment(teamId: string | number, departmentId: string | number) {
  return await req.delete(`api/teams/${teamId}/departments/${departmentId}`)
}

async function removeDepartmentMember(teamId: string | number, userId: string | number) {
  return await req.delete(`api/teams/${teamId}/departments/all/members/${userId}`)
}

