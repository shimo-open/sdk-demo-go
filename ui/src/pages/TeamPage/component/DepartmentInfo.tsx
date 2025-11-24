import { useModel } from 'umi'
import { useEffect, useState } from "react";
import { Tree, Card, Divider, Space, Button, Popconfirm, message } from "antd";
import { UsergroupAddOutlined, UserDeleteOutlined, PlusOutlined, CloseOutlined, GoldFilled, UserOutlined } from '@ant-design/icons';
import '../Team.module.less'
import { teamService } from '@/services/team.service'
import { CreateDepartmentModal } from './CreateDepartmentModal'
import { AddDepartmentMemberModal } from './AddDepartmentMemberModal'

export function DepartmentInfo() {
  const { userInfo } = useModel('user');
  const { teams } = useModel('team');
  const { currentUserId } = {
    currentUserId: userInfo?.id
  }
  const currentTeamId = teams && teams.length > 0 ? teams[0].id : 0
  const teamCreatorId = teams && teams.length > 0 ? teams[0].creator.id : 0
  console.log(teams);

  // Determine whether the current user is the team owner
  const isTeamCreator = Number(currentUserId) === Number(teamCreatorId)

  const [loaded, setLoaded] = useState(false)
  // Department tree data
  const [treeData, setTreeData] = useState([])
  const [selectedKey, setSelectedKey] = useState('')
  // Currently selected department
  const [selectedDept, setSelectedDept] = useState({} as any)
  // Currently selected member
  const [selectedUser, setSelectedUser] = useState({} as any)
  // Controls the add-department modal
  const [showAddDepartment, setShowAddDepartment] = useState(false)
  // Controls the add-member modal
  const [showAddDepartmentMember, setShowAddDepartmentMember] = useState(false)

  useEffect(() => {
    if (!loaded) {
      console.log('DepartmentInfo reloadTreeData')
      reloadTreeData()
    }
  }, [])

  function reloadTreeData() {
    teamService.departmentTopTree(currentTeamId).then((nodes) => {
      if (nodes.status === 200) {
        const newNodes = buildTree(nodes.data)
        setTreeData(newNodes)
        console.log('DepartmentInfo reloadTreeData', newNodes);
        setLoaded(true)
      }
    })
  }
  function buildTree(nodes: any) {
    // Derive a tree node ID based on the node type
    const getKeyByIdAndType = (id: string, type: string) => {
      switch (type) {
        case 'team':
          return 'root'
        case 'department':
          return `d-${id}`
        case 'user':
          return `u-${id}`
        default:
          return String(id)
      }
    }

    // Choose an icon based on the node type
    const getIconByType = (type: string) => {
      switch (type) {
        case 'team':
        case 'department':
          return <GoldFilled />
        case 'user':
          return <UserOutlined />
        default:
          return null
      }
    }

    return (nodes || []).map((n: any) => {
      const { node, type, children } = n
      const newNode: any = {
        title: node.name,
        key: getKeyByIdAndType(node.id, type),
        icon: getIconByType(type),
        isLeaf:
          type === 'user' ||
          (type === 'department' && (!children || children.length === 0))
      }

      if (children && children.length > 0) {
        newNode.children = buildTree(children)
      }

      return newNode
    })
  }

  const updateTreeData = (list: any, key: string, children: any) =>
    list.map((node: any) => {
      console.log('list', list, 'key', key, 'children', children)
      if (node.key === key) {
        return { ...node, children }
      }
      if (node.children) {
        return {
          ...node,
          children: updateTreeData(node.children, key, children)
        }
      }
      return node
    })

  // Delete a department
  function removeDepartment(id: string) {
    teamService.removeDepartment(currentTeamId, id).then(res => {
      reloadTreeData()
      setSelectedKey('')
      setSelectedDept({})
      message.success('删除成功')
    })
  }
  // Remove a member from a department
  function removeDepartmentMember(teamId: string, userId: string) {
    teamService.removeDepartmentMember(currentTeamId, userId).then(res => {
      reloadTreeData()
      setSelectedKey('')
      setSelectedUser({})
      message.success('移除成功')
    })
  }
  const onSelect = (keys: any, info: { event: "select"; selected: boolean; node: never; selectedNodes: never[]; nativeEvent: MouseEvent; }) => {
    console.log('Trigger Select', keys, info)
    if (keys.length === 1) {
      setSelectedKey(keys[0])
      // Department nodes
      if (String(keys[0]).startsWith('d-') || keys[0] === 'root') {
        setSelectedDept(info.node)
      }

      if (String(keys[0]).startsWith('u-')) {
        setSelectedUser(info.node)
      }
    } else if (keys.length === 0) {
      setSelectedKey('')
      setSelectedDept({})
      setSelectedUser({})
    }
  }
  const onLoadData = ({ key, children }: any) => {
    return new Promise((resolve) => {
      // Skip loading if children already exist
      if (children) {
        resolve(true)
        return
      }
      let departmentId = key

      if (String(departmentId).startsWith('d-')) {
        departmentId = String(departmentId).replace('d-', '')
      }

      teamService.getDepartmentChildren(currentTeamId, departmentId).then((nodes) => {
        console.log('load by key', key, 'children tree nodes ', nodes)
        const childrenNodes = buildTree(nodes.data)
        console.log('build children nodes ', childrenNodes)

        setTreeData((origin) => updateTreeData(origin, key, childrenNodes))
        resolve(true)
      })
    })
  }

  return (
    <Card>
      <h3>部门列表</h3>
      <Space style={{ marginBottom: '20px' }}>
        <Button type="primary" icon={<PlusOutlined />} disabled={
          !isTeamCreator ||
          !(selectedKey.indexOf('d') === 0 || selectedKey === 'root')
        } onClick={() => setShowAddDepartment(true)}>
          添加子部门
        </Button>
        <CreateDepartmentModal
          departmentId={(selectedKey && selectedKey.startsWith('d-')) ? selectedKey.replace('d-', '') : ''}
          onCreated={() => reloadTreeData()}
          open={showAddDepartment}
          onClose={() => setShowAddDepartment(false)}
        />
        <Button type="primary" icon={<UsergroupAddOutlined />}
          disabled={!isTeamCreator || !(selectedKey.indexOf('d') === 0)}
          onClick={() => setShowAddDepartmentMember(true)}>
          添加部门成员
        </Button>
        <AddDepartmentMemberModal
          departmentId={selectedDept && selectedDept.key ? selectedDept.key.replace('d-', '') : ''}
          onCreated={() => { reloadTreeData() }}
          open={showAddDepartmentMember}
          onClose={() => setShowAddDepartmentMember(false)}
        />
        <Popconfirm
          title="确认删除"
          description={`确认删除部门 「${selectedDept.title}」 ? 部门及所有下级部门和成员将被删除`}
          onConfirm={() => removeDepartment(String(selectedDept.key).replace('d-', ''))}
          okText="确认"
          cancelText="再想想"
        >

          <Button danger icon={<CloseOutlined />} disabled={!isTeamCreator || !(selectedKey.indexOf('d') === 0)}>删除部门</Button>
        </Popconfirm>
        <Popconfirm
          title="确认移除"
          description={`确认从部门中移除成员 「${selectedUser.title}」 ? 用户将不会从团队中移除`}
          onConfirm={() => removeDepartmentMember(currentTeamId, String(selectedUser.key).replace('u-', ''))}
          okText="确认"
          cancelText="再想想"
        >
          <Button danger icon={<UserDeleteOutlined />} disabled={!isTeamCreator || !(selectedKey.indexOf('u') === 0)}>移除部门成员</Button>
        </Popconfirm>
      </Space>
      <Tree
        showIcon
        showLine
        defaultExpandedKeys={['root']}
        onSelect={onSelect}
        loadData={onLoadData}
        treeData={treeData}
      />
      <Divider />
    </Card>
  )
}