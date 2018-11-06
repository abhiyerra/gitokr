import React from "react"

class Index extends React.Component {
    constructor(props) {
        super(props)

        this.state = {
            repo: "",
            tasks: [],
        }

        this.getTasks = this.getTasks.bind(this);
    }

    getTasks() {
        fetch(`http://localhost:8000/${this.state.repo}`).then(function (x) {
            console.log(x.json())
            this.setState({ tasks: x.json() })
        }.bind(this))

        return this.state.repo;
    }

    render() {
        return (
            <div>
                Repo:
                <select onChange={(x) => this.setState({ repo: x.target.value })}>
                    <option></option>
                    <option>acksin/consulting</option>
                </select>
                {this.getTasks()}

                {this.state.tasks}
            </div>
        )
    }
}

export default Index
