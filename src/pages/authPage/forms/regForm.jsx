import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Form, Input, Button, message } from "antd";
import classes from "../auth.module.css";

const RegForm = () => {
    const [form] = Form.useForm();
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState('');
    const navigate = useNavigate();

    // отправка нового юзера на сервер
    const sendNewUser = async (userData) => {
        try {
            const response = await fetch('http://localhost:8080/v1/auth/sign-up', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(userData),
                credentials: "include"
            });

            if (response.ok) {
                const result = await response.json();
                console.log('Успешная регистрация:', result);
                message.success('Регистрация прошла успешно!');
                setError('');

                setTimeout(() => {
                    navigate('/');
                }, 1000);

                return result;
            } else {
                const errorData = await response.json();
                console.error('Ошибка сервера:', response.status, errorData);

                const errorMessage = errorData.error || response.statusText;
                setError(errorMessage);
                message.error(`Ошибка регистрации: ${errorMessage}`);

                throw new Error(errorMessage);
            }
        } catch (error) {
            console.error('Ошибка:', error);
            const errorMsg = error.message || 'Произошла ошибка при регистрации';
            setError(errorMsg);
            message.error(errorMsg);
            throw error;
        }
    };

    // обработка отправки формы
    const onFinish = async (values) => {
        setLoading(true);
        setError('');
        try {
            const newUser = {
                user_name: values.userName,
                user_pass: values.userPass,
                user_email: values.userEmail
            };

            await sendNewUser(newUser);
            form.resetFields();

        } catch (error) {
        } finally {
            setLoading(false);
        }
    };

    const onFinishFailed = (errorInfo) => {
        console.log('Failed:', errorInfo);
        message.warning('Пожалуйста, заполните все поля правильно');
    };

    return (
        <div className={classes.formContainer}>
            {error && (
                <div className={classes.errorMessage} style={{
                    color: 'red',
                    textAlign: 'center',
                    marginBottom: '16px',
                    padding: '8px',
                    backgroundColor: '#fff2f0',
                    border: '1px solid #ffccc7',
                    borderRadius: '4px'
                }}>
                    {error}
                </div>
            )}
            <Form
                form={form}
                name="register"
                labelCol={{ span: 10 }}
                wrapperCol={{ span: 17 }}
                style={{ maxWidth: 600 }}
                onFinish={onFinish}
                onFinishFailed={onFinishFailed}
                autoComplete="off"
                validateTrigger="onSubmit"
            >
                <Form.Item
                    label="Name"
                    name="userName"

                    validateTrigger="onSubmit"
                    rules={[
                        {
                            required: true,
                            message: 'Please, enter your name!'
                        },
                        {
                            min: 3,
                            message: 'The name must be at least 3 characters long!'
                        }
                    ]}
                >
                    <Input />
                </Form.Item>

                <Form.Item
                    label="Email"
                    name="userEmail"
                    validateTrigger="onSubmit"
                    rules={[
                        {
                            required: true,
                            message: 'Please, enter your email!'
                        },
                        {
                            type: 'email',
                            message: 'Please enter a valid email address!'
                        }
                    ]}
                >
                    <Input />
                </Form.Item>

                <Form.Item
                    label="Password"
                    name="userPass"
                    validateTrigger="onSubmit"
                    rules={[
                        {
                            required: true,
                            message: 'Please, enter your password!'
                        },
                        {
                            min: 6,
                            message: 'The password must be at least 6 characters long!'
                        }
                    ]}
                >
                    <Input.Password />
                </Form.Item>
                <Form.Item
                    label="Confirm Password"
                    name="confirmPassword"
                    dependencies={['userPass']}
                    validateTrigger="onSubmit"
                    rules={[
                        {
                            required: true,
                            message: 'Please confirm your password!',
                        },
                        ({ getFieldValue }) => ({
                            validator(_, value) {
                                if (!value || getFieldValue('userPass') === value) {
                                    return Promise.resolve();
                                }
                                return Promise.reject(new Error('Passwords do not match!'));
                            },
                        }),
                    ]}
                >
                    <Input.Password />
                </Form.Item>

                <Form.Item wrapperCol={{ offset: 8, span: 16 }}>
                    <Button
                        type="primary"
                        htmlType="submit"
                        loading={loading}
                        className={classes.submitButton}
                    >
                        Sign UP!
                    </Button>
                </Form.Item>
            </Form>
        </div>
    );
}

export default RegForm;