<html>
    <head>
    <title></title>
    </head>
    <body>
        <table>
            <thead>
                <tr>
                   <td>{{.MinerID}}</td>
                   <td>Now</td> 
                </tr>
            </thead>
            <tbody>
                <tr>
                    <td>Available:</td>
                    <td>{{.Available}}</td>
                </tr>
                <tr>
                    <td>Pledged (FIL staked against stored data):</td>
                    <td>{{.Pledged}}</td>
                </tr>
                <tr>
                    <td>Locked (Rewards that haven't vested):</td>
                    <td>{{.Locked}}</td>
                </tr>
            </tbody>
        </table>
        <br></br>
        <table>
            <thead>
                <tr>
                    <td>Start Date: {{.StartDate}}</td>
                    <td></td>
                    <td>End Date: {{.EndDate}}</td>
                </tr>
                <tr>
                    <td>Assets</td>
                    <td>Cost</td>
                    <td>Revenue</td>
                </tr>
            </thead>
            <tbody>
                <tr>
                    <td>Transferred (Total FIL sent to the miner): {{.Transferred}}</td>
                    <td>Miner Fee (Tip fee given to miner that includes your message in a block): {{.MinerFee}}</td>
                    <td>Blocks won: {{.BlocksWon}}</td>
                </tr>
                <tr>
                    <td></td>
                    <td>Burn Fee (Total amount of burn fee across all message types): {{.BurnFee}}</td>
                    <td>FIL won: {{.FILWon}}</td>
                </tr>
                <tr>
                    <td></td>
                    <td>WindowPoSt (Burn fee for WindowPoSt messages): {{.WindowPoSt}}</td>
                    <td></td>
                </tr>
                <tr>
                    <td></td>
                    <td>PreCommit (Burn fee for PreCommit messages): {{.PreCommit}}</td>
                    <td></td>
                </tr>
                <tr>
                    <td></td>
                    <td>ProveCommit (Burn fee for ProveCommit messages): {{.ProveCommit}}</td>
                    <td></td>
                </tr>
                <tr>
                    <td></td>
                    <td>Penalty (FIL lost due to fault fees): {{.Penalty}}</td>
                    <td></td>
                </tr>
            </tbody>
        </table>
    </body>
</html>
