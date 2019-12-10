import React, { Component } from 'react';
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
  List,
  ListItem,
  ListInput,
  Row,
  Col,
  Gauge,
  Button
} from 'framework7-react';

export default class HomePage extends Component {
  state = {
    todayComplish: 0,
    todayTarget: 0,
    totalComplish: 0,
    totalTarget: 0,
  }

  componentDidMount() {
  }

  onPageAfterIn() {
    console.log("HomePage after in");
    this.checkLogin();
  }

  checkLogin() {
    reqwest({
      url: "./api/user",
      method: 'GET',
      type: 'json',
      error: (err) => {
        console.log("checkLogin error");
        this.$f7router.navigate('/login/');
      },
      success: (data) => {
        this.loadProgress();
      },
    });
  }

  loadProgress() {
    reqwest({
      url: "./api/progress",
      method: 'GET',
      type: 'json',
      error: (err) => {
      },
      success: (data) => {
        this.setState({
          todayComplish: data.tody,
          todayTarget: data.todayTarget,
          totalComplish: data.total,
          totalTarget: data.totalTarget,
        })
      },
    });
  }

  refreshProgress(e) {
    this.loadProgress();
  }

  startLearn(e) {
    this.$f7router.navigate('/learn/');
  }

  startReview(e) {
    this.$f7router.navigate('/reviewselection/');
  }

  render() {
    let todayComplish = this.state.todayTarget > 0 ? 1.0 * this.state.todayComplish / this.state.todayTarget : 0;
    let totalComplish = this.state.totalTarget > 0 ? 1.0 * this.state.totalComplish / this.state.totalTarget : 0;
    let todayTarget = "今日目标" + this.state.todayTarget;
    let totalLable = "总目标" + this.state.totalTarget;

    return (
      <Page
        onPageAfterIn={this.onPageAfterIn.bind(this)}>
        <Navbar>
          <NavLeft>
            <Link iconIos="f7:menu" iconMd="material:menu" panelOpen="left"></Link>
          </NavLeft>
          <NavTitle>背单词</NavTitle>
        </Navbar>
        <BlockTitle medium>进度</BlockTitle>
        <Block strong>
          <Row>
            <Col className="text-align-center">
              <Gauge
                type="semicircle"
                value={todayComplish}
                valueText={this.state.todayComplish}
                labelText={todayTarget}
                valueTextColor="#f44336"
                borderColor="#f44336"
              />
            </Col>
            <Col className="text-align-center">
              <Gauge
                type="semicircle"
                value={totalComplish}
                valueText={this.state.totalComplish}
                labelText={totalLable}
                valueTextColor="#e91e63"
                borderColor="#e91e63"
                labelTextColor="#333"
              />
            </Col>
          </Row>

          <Row>
            <Col width="50">
              <Button fill raised onClick={(e) => this.startLearn(e)}>继续学习</Button>
            </Col>
            <Col width="50">
              <Button fill raised onClick={(e) => this.startReview(e)}>开始复习</Button>
            </Col>
          </Row>
        </Block>

      </Page>
    );
  }
}
