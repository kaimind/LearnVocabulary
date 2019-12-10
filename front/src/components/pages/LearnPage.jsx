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
  BlockHeader,
  List,
  ListItem,
  Row,
  Col,
  Button
} from 'framework7-react';

export default class LearnPage extends Component {

  state = {
    userid: 1,
    index: 0,
    showDetail: false,
    words: []
  }

  onPageAfterIn() {
    this.loadWords();
  }

  loadWords() {
    reqwest({
      url: "./api/learn",
      method: 'GET',
      type: 'json',
      error: (err) => {
        this.$f7.preloader.hide();
        this.setState({
          index: 0,
          words: []
        });
      },
      success: (data) => {
        this.$f7.preloader.hide();
        this.setState({
          index: 0,
          words: data
        });
      },
    });
    this.$f7.preloader.show();
  }

  showWordDetail(e) {
    let showDetail = !this.state.showDetail;
    this.setState({
      showDetail: showDetail,
    });
  }

  showNextWord(e) {
    if (this.state.index < this.state.words.length - 1) {
      let index = this.state.index + 1;
      this.setState({
        showDetail: false,
        index: index,
      });
    } else if (this.state.words.length > 0) {
      this.putWordRecord();
    }
  }

  putWordRecord() {
    if (this.state.words.length <= 0) {
      return;
    }

    let words = []
    for (let i = 0; i < this.state.words.length; i++) {
      words.push(this.state.words[i].wid);
    }

    reqwest({
      url: "./api/learn",
      method: 'PUT',
      type: 'json',
      contentType: 'application/json',
      data: JSON.stringify({
        words: words
      }),
      error: (err) => {
        this.$f7.preloader.hide();
      },
      success: (data) => {
        // console.log("putwords, success!");
        // this.loadWords();

        // this.setState({
        //   index: 0,
        //   words: []
        // });
        this.$f7.preloader.hide();
        this.$f7router.back();
      }
    });
    this.$f7.preloader.show();
  }

  showPreWord(e) {
    if (this.state.index > 0) {
      let index = this.state.index - 1;
      this.setState({
        showDetail: false,
        index: index,
      });
    }
  }

  render() {
    let word = "无词可学";
    let read = []
    let meaning = []
    let detail = false
    let titleText = "学词 "

    if (this.state.words.length > 0) {
      titleText += "(" + (this.state.index + 1) + "/" + this.state.words.length + ")"

      let wordInfo = this.state.words[this.state.index]
      word = wordInfo.word
      for (let i = 0; wordInfo.detail.read && i < wordInfo.detail.read.length; i++) {
        read.push(<p key={i}>{wordInfo.detail.read[i].spell}</p>)
      }

      detail = this.state.showDetail
      if (detail && wordInfo.detail.define) {
        let keys = ['zh', 'form', 'dual', 'en']
        for (let i = 0; i < keys.length; i++) {
          if (wordInfo.detail.define[keys[i]]) {
            meaning.push(<p key={keys[i]}>{wordInfo.detail.define[keys[i]]}</p>)
          }
        }
      }
    }

    return (
      <Page
        onPageAfterIn={this.onPageAfterIn.bind(this)}>

        <Navbar title={titleText} backLink="返回" />

        <Toolbar bottom>
          <Link onClick={(e) => this.showPreWord(e)}>上一个</Link>
          <Link onClick={(e) => this.showWordDetail(e)}>不认识</Link>
          <Link onClick={(e) => this.showNextWord(e)}>下一个</Link>
        </Toolbar>

        <BlockTitle medium>{word}</BlockTitle>
        <Block strong>{read}</Block>

        {detail ? <BlockHeader>什么意思</BlockHeader> : null}
        {detail ? <Block strong>{meaning}</Block> : null}

      </Page>
    )
  }
}
