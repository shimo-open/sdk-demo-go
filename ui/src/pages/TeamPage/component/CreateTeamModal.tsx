import { useModel } from 'umi'
import { Modal, Form, Input, Button } from "antd";
import styles from '../Team.module.less'
import { teamContants } from '@/constants'
import { validate } from '@/utils/validate'

export function CreateTeamModal({ open, onClose }: { open: boolean, onClose: () => void }) {
  const [form] = Form.useForm();
  const { createTeam } = useModel('team');

  // Handle team creation
  function handleCreateTeam() {
    form.validateFields().then(() => {
      const formValue = form.getFieldsValue(true);
      createTeam(formValue.name)
      onClose()
      form.setFieldValue('name', '')
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
      title="创建团队"
      footer={[<Button onClick={onClose} key='cancel'>取消</Button>, <Button onClick={handleCreateTeam} key='create'>创建</Button>]}
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
          label="团队名称"
          name="name"
          rules={[{ required: true, message: '团队名称 1-20 字' }, { validator: validateItem('teamName') }]}
        >
          <Input placeholder="您的团队名称" />
        </Form.Item>
      </Form>
    </Modal>
  )
}
