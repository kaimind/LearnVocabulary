import React from 'react';
import { Page, Navbar, Block, BlockTitle, List, ListItem } from 'framework7-react';

export default class PannelLeftPage extends React.Component {
  render() {
    return (
      <Page>
        <Navbar title="关于" />
        <Block strong>
          <p>背单词</p>
        </Block>
        <BlockTitle>菜单</BlockTitle>
        <List>
          <ListItem link="/login/" title="退出" view="#main-view" panelClose></ListItem>
        </List>
      </Page>
    );
  }
}
