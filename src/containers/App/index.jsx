import React, { Component } from "react";
import getConfig from "services/Config";
import { init as initAlgolia } from "services/Algolia";
import Header from "components/Header";
import Herodote from "containers/Herodote";

/**
 * App Component.
 */
export default class App extends Component {
  /**
   * Creates an instance of App.
   * @param {Object} props Component props
   */
  constructor(props) {
    super(props);

    this.state = {};
  }

  /**
   * React lifecycle.
   */
  async componentDidMount() {
    const config = await getConfig();

    initAlgolia(config);

    this.setState({ config });
  }

  /**
   * React lifecycle.
   */
  render() {
    return (
      <div className="content">
        <Header />

        {this.state.config && <Herodote />}
      </div>
    );
  }
}
