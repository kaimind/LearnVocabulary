import React, { Component } from 'react';
import reqwest from 'reqwest';
import {
    Page,
    List,
    ListInput,
    ListButton,
    BlockFooter,
    LoginScreenTitle,
} from 'framework7-react';

class LoginPage extends Component {
    state = {
        username: '',
        password: '',
    }

    // componentDidMount() {
    //     this.onLogout();
    // }

    // onLogout() {
    //     reqwest({
    //         url: "/api/user",
    //         method: 'DELETE',
    //         type: 'json',
    //     });
    // }

    onLogin(e) {
        reqwest({
            url: "./api/user?username=" + this.state.username + "&password=" + this.state.password,
            method: 'POST',
            type: 'json',
            error: (err) => {
                this.$f7.preloader.hide();
                this.setState({
                    username: '',
                    password: '',
                });
            },
            success: (data) => {
                this.$f7.preloader.hide();
                this.setState({
                    username: '',
                    password: '',
                });

                this.$f7router.back();
            },
        });
        this.$f7.preloader.show();
    }

    render() {
        return (
            <Page noToolbar noNavbar noSwipeback loginScreen>
                <LoginScreenTitle>登录</LoginScreenTitle>
                <List form>
                    <ListInput
                        label="用户名"
                        type="text"
                        placeholder="用户名"
                        value={this.state.username}
                        onInput={(e) => {
                            this.setState({ username: e.target.value });
                        }}
                    />
                    <ListInput
                        label="密码"
                        type="password"
                        placeholder="密码"
                        value={this.state.password}
                        onInput={(e) => {
                            this.setState({ password: e.target.value });
                        }}
                    />
                </List>
                <List>
                    <ListButton onClick={(e) => this.onLogin(e)}>登录</ListButton>
                    <BlockFooter></BlockFooter>
                </List>
            </Page>
        )
    }
}

export default LoginPage;