<html>
    <head>
    <title>{{.Owner.MinerID}} Account Summary</title>
    </head>
    <body>
        <h1>Owner</h1>
        <h4>{{.Owner.MinerID}}</h4>
        <table>
            <tbody>
                <tr>
                    <td>Available:</td>
                    <td>{{.Owner.Available}}</td>
                </tr>
                <tr>
                    <td>Pledged:</td>
                    <td>{{.Owner.Pledged}}</td>
                </tr>
                <tr>
                    <td>Locked:</td>
                    <td>{{.Owner.Locked}}</td>
                </tr>
            </tbody>
        </table>
        <hr>
        <table>
            <thead>
                <tr>
                    <td>Start Date: {{.Owner.StartDate}}</td>
                    <td></td>
                    <td>End Date: {{.Owner.EndDate}}</td>
                </tr>
                <tr>
                    <td><b>Assets</b></td>
                    <td><b>Cost</b></td>
                    <td><b>Revenue</b></td>
                </tr>
            </thead>
            <tbody>
                <tr>
                    <td>Transferred: {{.Owner.Transferred}}</td>
                    <td>Miner Fee: {{.Owner.MinerFee}}</td>
                    <td>Blocks won: {{.Owner.BlocksWon}}</td>
                </tr>
                <tr>
                    <td></td>
                    <td>Burn Fee: {{.Owner.BurnFee}}</td>
                    <td>FIL won: {{.Owner.FILWon}}</td>
                </tr>
                <tr>
                    <td></td>
                    <td>WindowPoSt: {{.Owner.WindowPoSt}}</td>
                    <td></td>
                </tr>
                <tr>
                    <td></td>
                    <td>PreCommit: {{.Owner.PreCommit}}</td>
                    <td></td>
                </tr>
                <tr>
                    <td></td>
                    <td>ProveCommit: {{.Owner.ProveCommit}}</td>
                    <td></td>
                </tr>
                <tr>
                    <td></td>
                    <td>Penalty: {{.Owner.Penalty}}</td>
                    <td></td>
                </tr>
            </tbody>
        </table>
        <hr>
        <h1>Miners</h1>
        {{range .Miners}}
        <h4>{{.MinerID}}</h4>
        <table>
            <tbody>
                <tr>
                    <td>Available:</td>
                    <td>{{.Available}}</td>
                </tr>
                <tr>
                    <td>Pledged:</td>
                    <td>{{.Pledged}}</td>
                </tr>
                <tr>
                    <td>Locked:</td>
                    <td>{{.Locked}}</td>
                </tr>
            </tbody>
        </table>
        <hr>
        <table>
            <thead>
                <tr>
                    <td>Start Date: {{.StartDate}}</td>
                    <td></td>
                    <td>End Date: {{.EndDate}}</td>
                </tr>
                <tr>
                    <td><b>Assets</b></td>
                    <td><b>Cost</b></td>
                    <td><b>Revenue</b></td>
                </tr>
            </thead>
            <tbody>
                <tr>
                    <td>Transferred: {{.Transferred}}</td>
                    <td>Miner Fee: {{.MinerFee}}</td>
                    <td>Blocks won: {{.BlocksWon}}</td>
                </tr>
                <tr>
                    <td></td>
                    <td>Burn Fee: {{.BurnFee}}</td>
                    <td>FIL won: {{.FILWon}}</td>
                </tr>
                <tr>
                    <td></td>
                    <td>WindowPoSt: {{.WindowPoSt}}</td>
                    <td></td>
                </tr>
                <tr>
                    <td></td>
                    <td>PreCommit: {{.PreCommit}}</td>
                    <td></td>
                </tr>
                <tr>
                    <td></td>
                    <td>ProveCommit: {{.ProveCommit}}</td>
                    <td></td>
                </tr>
                <tr>
                    <td></td>
                    <td>Penalty: {{.Penalty}}</td>
                    <td></td>
                </tr>
            </tbody>
        </table>
        <hr>
        {{end}}
        <h3>Glossary</h3>
        <table>
            <tbody>
                <tr>
                    <td>Available:</td>
                    <td>Funds available to withdraw</td>
                </tr>
                <tr>
                    <td>Pledged:</td>
                    <td>FIL staked against stored data</td>
                </tr>
                <tr>
                    <td>Locked:</td>
                    <td>Rewards that haven't vested</td>
                </tr>
                <tr>
                    <td>Transferred:</td>
                    <td>Total FIL sent to the miner</td>
                </tr>
                <tr>
                    <td>Miner Fee:</td>
                    <td>Tip fee given to miner that includes your message in a block</td>
                </tr>
                <tr>
                    <td>Burn Fee:</td>
                    <td>Total amount of burn fee across all message types</td>
                </tr>
                <tr>
                    <td>WindowPoSt:</td>
                    <td>Burn fee for WindowPoSt messages</td>
                </tr>
                <tr>
                    <td>PreCommit:</td>
                    <td>Burn fee for PreCommit messages</td>
                </tr>
                <tr>
                    <td>ProveCommit:</td>
                    <td>Burn fee for ProveCommit messages</td>
                </tr>
                <tr>
                    <td>Penalty:</td>
                    <td>FIL lost due to fault fees</td>
                </tr>
            </tbody>
        </table>
    </body>
</html>
