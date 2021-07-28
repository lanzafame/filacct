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
                    <td>Pledged:</td>
                    <td>{{.Pledged}}</td>
                </tr>
                <tr>
                    <td>Locked:</td>
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
            </tbody>
        </table>
    </body>
</html>
