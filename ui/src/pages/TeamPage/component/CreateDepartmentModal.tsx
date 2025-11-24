import { useModel } from 'umi'
import { Modal, Form, Input, Button, message } from "antd";
import styles from '../Team.module.less'
import { teamService } from '@/services/team.service'
import { validate } from '@/utils/validate'

export function CreateDepartmentModal({ open, departmentId, onClose, onCreated }: { open: boolean, departmentId: number | string, onClose: () => void, onCreated: () => void }) {
  const [form] = Form.useForm();
  const { teams } = useModel('team');
  const currentTeamId = teams && teams.length > 0 ? teams[0].id : null

  // Handle creation
  function createTeam() {
    form.validateFields().then(() => {
      const formValue = form.getFieldsValue(true);
      teamService.addDepartment({
        teamId: currentTeamId,
        departmentId: departmentId,
        name: formValue.name
      }).then(res => {
        message.success('添加部门成功');
        form.resetFields();
      })
      onCreated()
      onClose()
    })
  }
  // Custom validation rule helper
  const validateItem = (name: keyof typeof validate) => (rule: any, value: string) => {
    return validate[name](value) as Promise<void>
  };
  type FieldType = { name: string }

  return (
    <Modal
      open={open}
      title="创建部门"
      footer={[<Button onClick={onClose} key='cancel'>取消</Button>, <Button onClick={createTeam} key='create'>创建</Button>]}
      onCancel={onClose}
      width="80vh"
      className={styles.modal}
    >
      <Form
        name="createTeam"
        form={form}
        labelCol={{ span: 5 }}
        wrapperCol={{ span: 18 }}
        autoComplete="off"
      >
        <Form.Item<FieldType>
          label="部门名称"
          name="name"
          rules={[{ required: true, message: '部门名称 1-30 字' }, { validator: validateItem('departmentName') }]}
        >
          <Input placeholder="新的部门名称" />
        </Form.Item>
      </Form>
    </Modal>
  )
}
