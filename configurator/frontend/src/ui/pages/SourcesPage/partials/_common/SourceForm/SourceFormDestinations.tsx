// @Libs
import React, { useCallback, useState } from 'react';
import { Link } from 'react-router-dom';
import { Avatar, Form, List, Switch } from 'antd';
// @Services
import ApplicationServices from '@service/ApplicationServices';
// @Components
import { getIconSrc } from '@page/DestinationsPage/partials/DestinationsList/DestinationsList';
import { loadDestinations } from '@page/DestinationsPage/commons';
import { CenteredError, CenteredSpin } from '@./lib/components/components';
// @Types
import { DestinationConfig } from '@service/destinations';
import { SourceFormDestinationsProps } from './SourceForm.types';
// @Hooks
import useLoader from '@./lib/commons/useLoader';
// @Routes
import { destinationPageRoutes } from '@page/DestinationsPage/DestinationsPage.routes';

const SourceFormDestinations = ({ initialValues, form }: SourceFormDestinationsProps) => {
  const [error, destinations] = useLoader(async() => await loadDestinations(ApplicationServices.get()));

  const [checkedDestinations, setCheckedDestinations] = useState<string[]>(initialValues.destinations ?? []);

  const handleChange = useCallback((config: DestinationConfig) => (checked: boolean) => {
    let newDestinations;

    if (checked && !checkedDestinations.includes(config.id)) {
      newDestinations = [...checkedDestinations, config.id];
    } else if (!checked) {
      newDestinations = checkedDestinations.filter((destination: string) => destination !== config.id);
    }

    setCheckedDestinations(newDestinations);

    const formValues = form.getFieldsValue();

    form.setFieldsValue({
      ...formValues,
      destinations: newDestinations
    });
  }, [checkedDestinations, form]);

  if (error) {
    return <CenteredError error={error} />
  } else if (!destinations) {
    return <CenteredSpin />
  }

  return (
    <>
      <h3>Choose destinations</h3>
      <article className="mb-5">
        <p>Destination is a database where reports data will be aggregated. Read more about destinations in our <a href="https://jitsu.com/docs/destinations-configuration" target="_blank" rel="noreferrer">documentation</a>.</p>
        {
          destinations.length > 0
            ? <>
              <p>You have to choose at least one destination.</p>
            </>
            : <p>If you haven't added any destinations yet you can do it <Link to={destinationPageRoutes.root}>here</Link>.</p>
        }
      </article>

      <Form.Item
        name="destinations"
        initialValue={initialValues.destinations}
        rules={destinations.length > 0 && [{ required: true, message: 'You have to choose at least one destination.' }]}
      >
        <List key="list" className="destinations-list" itemLayout="horizontal">
          {destinations.map((config) => {
            const description = config.describe();

            return <List.Item key={config.id}>
              <label htmlFor={config.id} className="ant-switch-group-label">
                <List.Item.Meta
                  avatar={<div className="ant-switch-group-label__avatar">
                    <Switch onChange={handleChange(config)} checked={checkedDestinations.includes(config.id)} />
                    <Avatar shape="square" src={getIconSrc(config.type)}/>
                  </div>}
                  description={<span className="destinations-list-show-connect-command">{description.displayURL}</span>}
                  title={config.connectionTestOk
                    ? config.id
                    : <span className="destinations-list-failed-connection">
                      <b>!</b> {config.id}
                    </span>}
                />
              </label>
            </List.Item>
          })}
        </List>
      </Form.Item>
    </>
  );
};

export { SourceFormDestinations }
