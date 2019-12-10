import React from 'react';
import reqwest from 'reqwest';
import {
    Page,
    Navbar,
    NavLeft,
    NavTitle,
    NavRight,
    Link,
    Toolbar,
    Block,
    BlockTitle,
    BlockHeader,
    List,
    ListItem,
    Row,
    Col,
    Button
} from 'framework7-react';

export default class ReviewSelection extends React.Component {
    state = {
        selected: [],
        learnList: []
    }

    onPageAfterIn() {
        console.log("ReviewSelectPage after in");
        this.loadList();
    }

    loadList() {

        reqwest({
            url: "./api/learnlist",
            method: 'GET',
            type: 'json',
            error: (err) => {
                this.$f7.preloader.hide();
            },
            success: (data) => {
                this.$f7.preloader.hide();
                this.setState({
                    selected: [],
                    learnList: data,
                })
            },
        });
        this.$f7.preloader.show();
    }

    resetCheckbox(e) {
        this.setState({
            selected: [],
        })
    }

    startReview(e) {
        if (this.state.selected.length <= 0) {
            return;
        }

        var learn = []
        for (let i = 0; i < this.state.selected.length; i++) {
            learn.push(this.state.selected[i]);
        }
        this.$f7router.navigate('/review/', {
            props: {
                learn: learn,
            }
        });
    }

    onCheckChange(e) {
        var self = this;
        var value = e.target.value;
        var selected = self.state.selected;
        if (e.target.checked) {
            selected.push(value);
        } else {
            selected.splice(selected.indexOf(value), 1);
        }
        self.setState({ selected: selected });
    }

    render() {
        var titleText = "选择 " + this.state.selected.length + " 天"
        var learnList = []
        for (let i = 0; i < this.state.learnList.length; i++) {
            let checked = this.state.selected.indexOf(this.state.learnList[i].learn) >= 0;
            learnList.push(<ListItem
                checkbox
                name="demo-media-checkbox"
                checked={checked}
                key={this.state.learnList[i].learn}
                value={this.state.learnList[i].learn}
                title={this.state.learnList[i].learn}
                after={this.state.learnList[i].review}
                onChange={(e) => this.onCheckChange(e)} />);
        }

        return (<Page
            onPageAfterIn={this.onPageAfterIn.bind(this)}>

            <Navbar title={titleText} backLink="返回" />

            <Toolbar bottom>
                <Link onClick={(e) => this.resetCheckbox(e)}>重置</Link>
                <Link onClick={(e) => this.startReview(e)}>开始</Link>
            </Toolbar>

            <List mediaList>
                {learnList}
            </List>

        </Page>)
    }
}