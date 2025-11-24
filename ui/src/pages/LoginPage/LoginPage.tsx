
import { history, useModel } from "umi"
import { otherConstants } from '@/constants/other.constants'
import styles from './LoginPage.module.less'
import { CopyrightFooter } from '../CopyrightFooter/Footer'
import { Tabs, Button, Form, Input, message } from "antd";
import type { TabsProps, FormProps } from 'antd';
import { validate } from '@/utils/validate'
import { userService } from '@/services/user.service';
import { setToken } from '@/utils/axios'
import { useState } from "react";
import { userConstants } from "@/constants/user.constants";

export default function LoginPage() {
    const items: TabsProps['items'] = [
        {
            key: 'login',
            label: '登录',
            children: <LoginOrRegister type="login" />
        },
        {
            key: 'register',
            label: '注册',
            children: <LoginOrRegister type="register" />
        }
    ]

    return (
        <div className={styles.container}>
            <img src={otherConstants.LOGO} className={styles.logo} alt="石墨文档" />
            <Tabs defaultActiveKey="login" items={items} />
            <CopyrightFooter className={styles.copyrightFooter} />
        </div>
    );
}

type FCProps = { type: string };

type FieldType = {
    email?: string;
    password?: string
};

function LoginOrRegister({ type }: FCProps) {
    const { setUserInfo, setHostPath } = useModel("user");
    const [loading, setLoading] = useState(false)
    const { setAvailableFileTypes, setPremiumFileTypes, setAllowBoardAndMindmap } = useModel("file")
    const onFinish: FormProps<FieldType>['onFinish'] = (values: any) => {
        setLoading(true)
        const { email, password } = values;
        if (type === 'login') {
            // Handle login
            userService.login({
                email: email,
                password: password
            }).then(res => {
                handleUserInfo(res.data)
            }).finally(() => {
                setLoading(false)
            })
        } else if (type === 'register') {
            // Handle registration
            userService.register({
                email: email,
                password: password
            }).then(res => {
                handleUserInfo(res.data)
            }).finally(() => {
                setLoading(false)
            })
        }
    };

    // Process user info after authentication
    const handleUserInfo = async ({ user, token, availableFileTypes, premiumFileTypes, hostPath, allowBoardAndMindmap }: any) => {
        localStorage.setItem(userConstants.LOGIN_KEY, userConstants.LOGIN_KEY)
        setUserInfo(user)
        setUserInfo(user)
        setToken(token);
        setAvailableFileTypes(availableFileTypes)
        setPremiumFileTypes(premiumFileTypes)
        setAllowBoardAndMindmap(allowBoardAndMindmap)
        setHostPath(hostPath)
        const redirect = new URLSearchParams(location.search)?.get(userConstants.REDIRECT_URL_KEY) || ""
        if (redirect !== "") {
            history.push(redirect);
        } else {
            history.push(`/home`);
        }
    }

    // Custom validation rules
    const validateItem = (name: keyof typeof validate) => (rule: any, value: string) => {
        return validate[name](value) as Promise<void>
    };

    return (
        <Form
            name={type}
            labelCol={{ span: 4 }}
            wrapperCol={{ span: 16 }}
            onFinish={onFinish}
            autoComplete="off"
        >
            <Form.Item<FieldType>
                label="邮箱"
                name="email"
                rules={[{ required: true, message: '请输入邮箱' }, { validator: validateItem('email') }]}
            >
                <Input placeholder="您的邮箱" />
            </Form.Item>

            <Form.Item<FieldType>
                label="密码"
                name="password"
                rules={[{ required: true, message: '请输入密码' }, { validator: validateItem('password') }]}
            >
                <Input.Password placeholder="您的密码" />
            </Form.Item>

            <Form.Item wrapperCol={{ offset: 7 }}>
                <Button type="primary" htmlType="submit" loading={loading}>
                    {type == 'login' ? '登录' : '注册'}
                </Button>
            </Form.Item>
        </Form>
    );
}
