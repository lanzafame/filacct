<!doctype html>
<html>
    <head>
      <meta charset="UTF-8" />
      <meta name="viewport" content="width=device-width, initial-scale=1.0" />
      <link href="/tailwind.css" rel="stylesheet">
    <title>{{.MinerID}} Account Summary</title>
    </head>
    <body>
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
                <tr>
                    <td></td>
                    <td>Sent out: {{.Sent}}</td>
                    <td></td>
                </tr>
            </tbody>
        </table>
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

